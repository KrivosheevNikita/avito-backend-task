package models

import "errors"

var (
	ErrBadRequest     = errors.New("Неверный запрос")
	ErrInternalServer = errors.New("Внутренняя ошибка сервера")
	ErrUnauthorized   = errors.New("Неверные учетные данные")
	ErrForbidden      = errors.New("Доступ запрещен")
)

var (
	ErrInvalidRole = errors.New("Неверная роль")
	ErrEmailExist  = errors.New("Пользователь с данной почтой уже зарегестрирован")
)

var (
	ErrPVZNotFound       = errors.New("ПВЗ не найден")
	ErrInvalidCity       = errors.New("City должен быть один из: Москва, Санкт-Петербург, Казань")
	ErrInvalidDateRange  = errors.New("Неверная дата")
	ErrInvalidPagination = errors.New("Неверный номер страницы или лимит")
	ErrPvzID             = errors.New("Неверный PvzId")
)

var (
	ErrNoOpenReception        = errors.New("Нет открытых приемок для закрытия")
	ErrReceptionAlreadyClosed = errors.New("Приемка уже закрыта")
	ErrReceptionInProgress    = errors.New("Уже есть активная приемка в статусе in_progress")
	ErrNoActiveReception      = errors.New("Нет активной приемки для добавления товара")
)

var (
	ErrNoProductToDelete  = errors.New("Нет товара для удаления")
	ErrInvalidTypeProduct = errors.New("Неправильный тип товара")
)
