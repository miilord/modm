package modm

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DoTransactionFunc is a function signature for performing transactions.
type DoTransactionFunc func(
	ctx context.Context,
	callback func(sessCtx context.Context) (interface{}, error),
	opts ...*options.TransactionOptions,
) (interface{}, error)

// DoTransaction creates and manages a database transaction.
func DoTransaction(client *mongo.Client) DoTransactionFunc {
	return func(
		ctx context.Context,
		callback func(sessCtx context.Context) (interface{}, error),
		opts ...*options.TransactionOptions,
	) (interface{}, error) {
		sess, err := client.StartSession()
		if err != nil {
			return nil, err
		}
		defer sess.EndSession(ctx)

		res, err := sess.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
			return callback(sessCtx)
		}, opts...)
		if err != nil {
			return nil, err
		}

		return res, nil
	}
}
