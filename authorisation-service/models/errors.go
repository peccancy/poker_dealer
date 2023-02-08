package models

type ErrNotValidIDToken string

func (e ErrNotValidIDToken) Error() string {
	return string(e)
}

type ErrUserDenied string

func (e ErrUserDenied) Error() string {
	return string(e)
}

type ErrInternalSrvErr struct{}

func (e ErrInternalSrvErr) Error() string {
	const msg = "internal service error"
	return msg
}
