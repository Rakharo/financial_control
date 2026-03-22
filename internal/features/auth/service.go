package auth

import (
	"encoding/json"
	"errors"
	"financial_control/internal/features/user"
	"financial_control/internal/utils"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	authRepo Repository
	userRepo user.Repository
}

func NewService(authRepo Repository, userRepo user.Repository) *Service {
	return &Service{authRepo: authRepo, userRepo: userRepo}
}

func (s *Service) UpdateUserPassword(userID uint64, password PasswordRequest) error {
	user, err := s.authRepo.GetUserWithPasswordById(userID)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("usuário não encontrado")
	}

	//validacao da senha atual
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password.CurrentPassword),
	)
	if err != nil {
		return errors.New("senha incorreta")
	}

	//validacao de senha nova e confirmacao de senha nova
	if password.NewPassword != password.ConfirmPassword {
		return errors.New("as senhas não coincidem")
	}

	//validacao de regra da senha
	if err := utils.ValidatePassword(password.NewPassword); err != nil {
		return err
	}

	//nova criptografia da senha
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password.NewPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	return s.authRepo.UpdateUserPassword(userID, string(hash))
}

func (s *Service) Login(email string, password string) (*user.User, string, string, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, "", "", utils.NewUnauthorized("E-mail ou senha inválidos!", nil)
	}

	if user == nil {
		return nil, "", "", utils.NewUnauthorized("E-mail ou senha inválidos!", nil)
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)
	if err != nil {
		return nil, "", "", utils.NewUnauthorized("E-mail ou senha inválidos", err)
	}

	accessToken, refreshToken := s.createTokens(user.ID)

	return user, accessToken, refreshToken, nil
}

func (s *Service) LoginWithGoogle(token string) (*user.User, string, string, error) {
	resp, err := http.Get("https://oauth2.googleapis.com/tokeninfo?id_token=" + token)
	if err != nil {
		return nil, "", "", utils.NewUnauthorized("Token inválido", err)
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, "", "", utils.NewUnauthorized("Erro ao validar token", err)
	}

	if data["aud"] != "886522038636-kbjmui0f7t0h4lcdg4lrfc7sik168jsu.apps.googleusercontent.com" {
		return nil, "", "", utils.NewUnauthorized("Token inválido para este app", nil)
	}

	email, _ := data["email"].(string)
	name, _ := data["name"].(string)
	providerID, _ := data["sub"].(string)

	if email == "" || providerID == "" {
		return nil, "", "", utils.NewUnauthorized("E-mail ou ID não encontrado no token", nil)
	}

	up, err := s.authRepo.GetUserProvider("google", providerID)
	if err != nil {
		return nil, "", "", err
	}

	var u *user.User

	if up != nil {
		// existe → pega usuário pelo user_id
		u, err = s.userRepo.GetUserById(uint64(up.UserID))
		if err != nil {
			return nil, "", "", err
		}
	} else {
		// não existe → cria usuário e user_provider
		u, err = s.userRepo.GetUserByEmail(email)
		if err != nil {
			return nil, "", "", err
		}

		if u == nil {
			// cria novo usuário
			newUser := &user.User{
				Name:  name,
				Email: email,
			}
			if err := s.userRepo.CreateUser(newUser); err != nil {
				return nil, "", "", err
			}
			u = newUser
		}

		// cria user_provider vinculado
		newUP := &UserProvider{
			UserID:         int(u.ID),
			ProviderName:   "google",
			ProviderUserID: providerID,
			CreatedAt:      time.Now(),
		}
		if err := s.authRepo.CreateUserProvider(newUP); err != nil {
			return nil, "", "", err
		}
	}

	accessToken, refreshToken := s.createTokens(u.ID)

	return u, accessToken, refreshToken, nil
}

func (s *Service) LinkGoogleAccount(userID uint64, token string) error {
	resp, err := http.Get("https://oauth2.googleapis.com/tokeninfo?id_token=" + token)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}

	email, _ := data["email"].(string)
	providerID, _ := data["sub"].(string)

	if email == "" || providerID == "" {
		return errors.New("dados inválidos do Google")
	}

	user, err := s.userRepo.GetUserById(userID)
	if err != nil || user == nil {
		return errors.New("usuário não encontrado")
	}

	if user.Email != email {
		return errors.New("email do Google diferente do usuário logado")
	}

	existing, err := s.authRepo.GetUserProvider("google", providerID)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("conta Google já vinculada")
	}

	up := &UserProvider{
		UserID:         int(userID),
		ProviderName:   "google",
		ProviderUserID: providerID,
		CreatedAt:      time.Now(),
	}

	return s.authRepo.CreateUserProvider(up)
}

func (s *Service) RefreshToken(token string) (string, string, error) {
	rt, err := s.authRepo.GetRefreshToken(token)
	if time.Now().After(rt.ExpiresAt) {
		_ = s.authRepo.DeleteRefreshToken(token)
		return "", "", utils.NewUnauthorized("Refresh token expirado", nil)
	}
	if err != nil {
		return "", "", err
	}
	if rt == nil {
		return "", "", utils.NewUnauthorized("Refresh token inválido", nil)
	}

	// gera novos tokens
	accessToken, refreshToken := s.createTokens(uint64(rt.UserID))

	// remove antigo refresh token
	_ = s.authRepo.DeleteRefreshToken(token)

	return accessToken, refreshToken, nil
}

func (s *Service) Logout(userID uint64, token string) error {
	rt, err := s.authRepo.GetRefreshToken(token)
	if err != nil {
		return err
	}

	if rt == nil || uint64(rt.UserID) != userID {
		return errors.New("token inválido")
	}

	return s.authRepo.DeleteRefreshToken(token)
}

func (s *Service) createTokens(userID uint64) (string, string) {
	accessToken, _ := GenerateAccessToken(userID)
	refreshToken, _ := GenerateRefreshToken(userID)

	// salva refresh token no banco
	rt := &RefreshToken{
		UserID:    int(userID),
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now(),
	}
	err := s.authRepo.CreateRefreshToken(rt)
	if err != nil {
		return "", ""
	}

	return accessToken, refreshToken
}

func (s *Service) GetUserProviders(userID uint64) ([]string, error) {
	return s.authRepo.GetProvidersByUserID(userID)
}
