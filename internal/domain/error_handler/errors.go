package handler

import "errors"

var (
	// User
	ErrEmptyName  = errors.New("nome não pode ser vazio")
	ErrEmptyEmail = errors.New("email não pode ser vazio")

	// Location
	ErrInvalidCPTECCode  = errors.New("código CPTEC inválido")
	ErrEmptyLocationName = errors.New("nome da localização não pode ser vazio")
	ErrInvalidState      = errors.New("estado deve ter 2 caracteres")

	// Notification
	ErrInvalidUserID             = errors.New("ID do usuário inválido")
	ErrInvalidLocationID         = errors.New("ID da localização inválido")
	ErrInvalidScheduleDate       = errors.New("data de agendamento deve ser futura")
	ErrInvalidNotificationStatus = errors.New("status da notificação inválido para envio")
	ErrEmptyForecast             = errors.New("previsão do tempo não pode estar vazia")

	// Service
	ErrUserOptOut          = errors.New("usuário optou por não receber notificações")
	ErrInvalidScheduleTime = errors.New("horário de agendamento inválido")
	ErrCPTECUnavailable    = errors.New("serviço CPTEC indisponível")

	// Repository
	ErrNotFound     = errors.New("registro não encontrado")
	ErrDuplicateKey = errors.New("chave duplicada")
	ErrInvalidInput = errors.New("entrada inválida")
)
