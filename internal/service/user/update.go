package note

import (
	"context"

	"auth/internal/model"
)

func (s *serv) Update(ctx context.Context, id int64, info *model.UserInfo) error {

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		errTx = s.userRepository.Update(ctx, id, info)
		if errTx != nil {
			return errTx
		}

		_, errTx = s.userRepository.Get(ctx, id)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}
