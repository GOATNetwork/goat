package keeper_test

func (suite *KeeperTestSuite) TestUpdateForkParams() {
	suite.Context = suite.Context.WithChainID("unitest").WithBlockHeight(10)
	err := suite.Keeper.UpdateForkParams(suite.Context)
	suite.Require().NoError(err)
	param, err := suite.Keeper.Params.Get(suite.Context)
	suite.Require().NoError(err)
	suite.Require().EqualValues(2186000000000000000, param.InitialBlockReward)
}
