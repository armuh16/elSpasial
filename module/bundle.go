package module

import (
	//Route
	authRoute "github.com/elspasial/module/auth/route"
	transactionRoute "github.com/elspasial/module/transaction/route"
	productRoute "github.com/elspasial/module/trip/route"

	//Logic
	authLogic "github.com/elspasial/module/auth/logic"
	transactionLogic "github.com/elspasial/module/transaction/logic"
	productLogic "github.com/elspasial/module/trip/logic"
	userLogic "github.com/elspasial/module/user/logic"

	//Repository
	transactionRepository "github.com/elspasial/module/transaction/repository"
	productRepository "github.com/elspasial/module/trip/repository"
	userRepository "github.com/elspasial/module/user/repository"

	"go.uber.org/fx"
)

// Register Route
var BundleRoute = fx.Options(
	fx.Invoke(transactionRoute.NewRoute),
	fx.Invoke(productRoute.NewRoute),
	fx.Invoke(authRoute.NewRoute),
)

// Register logic
var BundleLogic = fx.Options(
	fx.Provide(userLogic.NewLogic),
	fx.Provide(transactionLogic.NewLogic),
	fx.Provide(productLogic.NewLogic),
	fx.Provide(authLogic.NewLogic),
)

// Register Repository
var BundleRepository = fx.Options(
	fx.Provide(userRepository.NewRepository),
	fx.Provide(transactionRepository.NewRepository),
	fx.Provide(productRepository.NewRepository),
)
