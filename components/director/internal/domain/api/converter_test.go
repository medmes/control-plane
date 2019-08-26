package api_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/kyma-incubator/compass/components/director/pkg/strings"

	"github.com/kyma-incubator/compass/components/director/internal/repo/testdb"

	"github.com/google/uuid"
	"github.com/kyma-incubator/compass/components/director/internal/domain/version"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/kyma-incubator/compass/components/director/internal/model"

	"github.com/kyma-incubator/compass/components/director/internal/domain/api"
	"github.com/kyma-incubator/compass/components/director/internal/domain/api/automock"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
	"github.com/stretchr/testify/assert"
)

func TestConverter_ToGraphQL(t *testing.T) {
	// given
	modelAPIDefinition := fixDetailedModelAPIDefinition("foo", "Foo", "Lorem ipsum", "group")
	gqlAPIDefinition := fixDetailedGQLAPIDefinition("foo", "Foo", "Lorem ipsum", "group")
	emptyModelAPIDefinition := &model.APIDefinition{}
	emptyGraphQLAPIDefinition := &graphql.APIDefinition{}

	testCases := []struct {
		Name                  string
		Input                 *model.APIDefinition
		Expected              *graphql.APIDefinition
		AuthConverterFn       func() *automock.AuthConverter
		FetchRequestConverter func() *automock.FetchRequestConverter
		VersionConverter      func() *automock.VersionConverter
	}{
		{
			Name:     "All properties given",
			Input:    modelAPIDefinition,
			Expected: gqlAPIDefinition,
			AuthConverterFn: func() *automock.AuthConverter {
				conv := &automock.AuthConverter{}
				conv.On("ToGraphQL", modelAPIDefinition.DefaultAuth).Return(gqlAPIDefinition.DefaultAuth).Once()
				for i, auth := range modelAPIDefinition.Auths {
					conv.On("ToGraphQL", auth.Auth).Return(gqlAPIDefinition.Auths[i].Auth).Once()

				}

				return conv
			},
			FetchRequestConverter: func() *automock.FetchRequestConverter {
				conv := &automock.FetchRequestConverter{}
				conv.On("ToGraphQL", modelAPIDefinition.Spec.FetchRequest).Return(gqlAPIDefinition.Spec.FetchRequest).Once()
				return conv
			},
			VersionConverter: func() *automock.VersionConverter {
				conv := &automock.VersionConverter{}
				conv.On("ToGraphQL", modelAPIDefinition.Version).Return(gqlAPIDefinition.Version).Once()
				return conv
			},
		},
		{
			Name:     "Empty",
			Input:    emptyModelAPIDefinition,
			Expected: emptyGraphQLAPIDefinition,
			AuthConverterFn: func() *automock.AuthConverter {
				conv := &automock.AuthConverter{}
				conv.On("ToGraphQL", emptyModelAPIDefinition.DefaultAuth).Return(nil).Once()
				return conv
			},
			FetchRequestConverter: func() *automock.FetchRequestConverter {
				return &automock.FetchRequestConverter{}
			},
			VersionConverter: func() *automock.VersionConverter {
				conv := &automock.VersionConverter{}
				conv.On("ToGraphQL", emptyModelAPIDefinition.Version).Return(nil).Once()
				return conv
			},
		},
		{
			Name:     "Nil",
			Input:    nil,
			Expected: nil,
			AuthConverterFn: func() *automock.AuthConverter {
				return &automock.AuthConverter{}
			},
			FetchRequestConverter: func() *automock.FetchRequestConverter {
				return &automock.FetchRequestConverter{}
			},
			VersionConverter: func() *automock.VersionConverter {
				return &automock.VersionConverter{}
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s", testCase.Name), func(t *testing.T) {
			//given
			authConverter := testCase.AuthConverterFn()
			frConverter := testCase.FetchRequestConverter()
			versionConverter := testCase.VersionConverter()

			// when
			converter := api.NewConverter(authConverter, frConverter, versionConverter)
			res := converter.ToGraphQL(testCase.Input)

			// then
			assert.EqualValues(t, testCase.Expected, res)
			authConverter.AssertExpectations(t)
			frConverter.AssertExpectations(t)
			versionConverter.AssertExpectations(t)
		})
	}
}

func TestConverter_MultipleToGraphQL(t *testing.T) {
	// given
	input := []*model.APIDefinition{
		fixModelAPIDefinition("foo", "1", "Foo", "Lorem ipsum"),
		fixModelAPIDefinition("bar", "1", "Bar", "Dolor sit amet"),
		{},
		nil,
	}

	expected := []*graphql.APIDefinition{
		fixGQLAPIDefinition("foo", "1", "Foo", "Lorem ipsum"),
		fixGQLAPIDefinition("bar", "1", "Bar", "Dolor sit amet"),
		{},
	}

	authConverter := &automock.AuthConverter{}
	frConverter := &automock.FetchRequestConverter{}
	versionConverter := &automock.VersionConverter{}

	for i, api := range input {
		if api == nil {
			continue
		}
		authConverter.On("ToGraphQL", api.DefaultAuth).Return(expected[i].DefaultAuth).Once()
		versionConverter.On("ToGraphQL", api.Version).Return(expected[i].Version).Once()
	}

	// when
	converter := api.NewConverter(authConverter, frConverter, versionConverter)
	res := converter.MultipleToGraphQL(input)

	// then
	assert.Equal(t, expected, res)
	authConverter.AssertExpectations(t)
	frConverter.AssertExpectations(t)
}

func TestConverter_InputFromGraphQL(t *testing.T) {
	// given
	gqlAPIDefinitionInput := fixGQLAPIDefinitionInput("foo", "Lorem ipsum", "group")
	modelAPIDefinitionInput := fixModelAPIDefinitionInput("foo", "Lorem ipsum", "group")
	emptyGQLAPIDefinition := &graphql.APIDefinitionInput{}
	emptyModelAPIDefinition := &model.APIDefinitionInput{}
	testCases := []struct {
		Name                  string
		Input                 *graphql.APIDefinitionInput
		Expected              *model.APIDefinitionInput
		AuthConverterFn       func() *automock.AuthConverter
		FetchRequestConverter func() *automock.FetchRequestConverter
		VersionConverter      func() *automock.VersionConverter
	}{
		{
			Name:     "All properties given",
			Input:    gqlAPIDefinitionInput,
			Expected: modelAPIDefinitionInput,
			AuthConverterFn: func() *automock.AuthConverter {
				conv := &automock.AuthConverter{}
				conv.On("InputFromGraphQL", gqlAPIDefinitionInput.DefaultAuth).Return(modelAPIDefinitionInput.DefaultAuth).Once()
				return conv
			},
			FetchRequestConverter: func() *automock.FetchRequestConverter {
				conv := &automock.FetchRequestConverter{}
				conv.On("InputFromGraphQL", gqlAPIDefinitionInput.Spec.FetchRequest).Return(modelAPIDefinitionInput.Spec.FetchRequest).Once()
				return conv
			},
			VersionConverter: func() *automock.VersionConverter {
				conv := &automock.VersionConverter{}
				conv.On("InputFromGraphQL", gqlAPIDefinitionInput.Version).Return(modelAPIDefinitionInput.Version).Once()
				return conv
			},
		},
		{
			Name:     "Empty",
			Input:    &graphql.APIDefinitionInput{},
			Expected: &model.APIDefinitionInput{},
			AuthConverterFn: func() *automock.AuthConverter {
				conv := &automock.AuthConverter{}
				conv.On("InputFromGraphQL", emptyGQLAPIDefinition.DefaultAuth).Return(emptyModelAPIDefinition.DefaultAuth).Once()
				return conv
			},
			FetchRequestConverter: func() *automock.FetchRequestConverter {
				return &automock.FetchRequestConverter{}
			},
			VersionConverter: func() *automock.VersionConverter {
				conv := &automock.VersionConverter{}
				conv.On("InputFromGraphQL", emptyGQLAPIDefinition.Version).Return(nil).Once()
				return conv
			},
		},
		{
			Name:     "Nil",
			Input:    nil,
			Expected: nil,
			AuthConverterFn: func() *automock.AuthConverter {
				return &automock.AuthConverter{}
			},
			FetchRequestConverter: func() *automock.FetchRequestConverter {
				return &automock.FetchRequestConverter{}
			},
			VersionConverter: func() *automock.VersionConverter {
				return &automock.VersionConverter{}
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s", testCase.Name), func(t *testing.T) {
			//given
			authConverter := testCase.AuthConverterFn()
			frConverter := testCase.FetchRequestConverter()
			versionConverter := testCase.VersionConverter()

			// when
			converter := api.NewConverter(authConverter, frConverter, versionConverter)
			res := converter.InputFromGraphQL(testCase.Input)

			// then
			assert.Equal(t, testCase.Expected, res)
			authConverter.AssertExpectations(t)
			frConverter.AssertExpectations(t)
			versionConverter.AssertExpectations(t)
		})
	}
}

func TestConverter_MultipleInputFromGraphQL(t *testing.T) {
	// given
	gqlAPI1 := fixGQLAPIDefinitionInput("foo", "lorem", "group")
	gqlAPI2 := fixGQLAPIDefinitionInput("bar", "ipsum", "group2")

	modelAPI1 := fixModelAPIDefinitionInput("foo", "lorem", "group")
	modelAPI2 := fixModelAPIDefinitionInput("bar", "ipsum", "group2")

	gqlAPIDefinitionInputs := []*graphql.APIDefinitionInput{gqlAPI1, gqlAPI2}
	modelAPIDefinitionInputs := []*model.APIDefinitionInput{modelAPI1, modelAPI2}
	testCases := []struct {
		Name                  string
		Input                 []*graphql.APIDefinitionInput
		Expected              []*model.APIDefinitionInput
		AuthConverterFn       func() *automock.AuthConverter
		FetchRequestConverter func() *automock.FetchRequestConverter
		VersionConverter      func() *automock.VersionConverter
	}{
		{
			Name:     "All properties given",
			Input:    gqlAPIDefinitionInputs,
			Expected: modelAPIDefinitionInputs,
			AuthConverterFn: func() *automock.AuthConverter {
				conv := &automock.AuthConverter{}
				for i, apiDef := range gqlAPIDefinitionInputs {
					conv.On("InputFromGraphQL", apiDef.DefaultAuth).Return(modelAPIDefinitionInputs[i].DefaultAuth).Once()
				}

				return conv
			},
			FetchRequestConverter: func() *automock.FetchRequestConverter {
				conv := &automock.FetchRequestConverter{}
				for i, apiDef := range gqlAPIDefinitionInputs {
					conv.On("InputFromGraphQL", apiDef.Spec.FetchRequest).Return(modelAPIDefinitionInputs[i].Spec.FetchRequest).Once()

				}

				return conv
			},
			VersionConverter: func() *automock.VersionConverter {
				conv := &automock.VersionConverter{}
				for i, apiDef := range gqlAPIDefinitionInputs {
					conv.On("InputFromGraphQL", apiDef.Version).Return(modelAPIDefinitionInputs[i].Version).Once()
				}
				return conv
			},
		},
		{
			Name:     "Empty",
			Input:    []*graphql.APIDefinitionInput{},
			Expected: nil,
			AuthConverterFn: func() *automock.AuthConverter {
				return &automock.AuthConverter{}
			},
			FetchRequestConverter: func() *automock.FetchRequestConverter {
				return &automock.FetchRequestConverter{}
			},
			VersionConverter: func() *automock.VersionConverter {
				return &automock.VersionConverter{}
			},
		},
		{
			Name:     "Nil",
			Input:    nil,
			Expected: nil,
			AuthConverterFn: func() *automock.AuthConverter {
				return &automock.AuthConverter{}
			},
			FetchRequestConverter: func() *automock.FetchRequestConverter {
				return &automock.FetchRequestConverter{}
			},
			VersionConverter: func() *automock.VersionConverter {
				return &automock.VersionConverter{}
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s", testCase.Name), func(t *testing.T) {
			//given
			authConverter := testCase.AuthConverterFn()
			frConverter := testCase.FetchRequestConverter()
			versionConverter := testCase.VersionConverter()

			// when
			converter := api.NewConverter(authConverter, frConverter, versionConverter)
			res := converter.MultipleInputFromGraphQL(testCase.Input)

			// then
			assert.Equal(t, testCase.Expected, res)
			authConverter.AssertExpectations(t)
			frConverter.AssertExpectations(t)
			versionConverter.AssertExpectations(t)
		})
	}
}

func TestApiSpecDataConversionNilStaysNil(t *testing.T) {
	// GIVEN
	mockAuthConv := &automock.AuthConverter{}
	defer mockAuthConv.AssertExpectations(t)
	mockAuthConv.On("InputFromGraphQL", mock.Anything).Return(nil)
	mockAuthConv.On("ToGraphQL", mock.Anything).Return(nil)

	mockFrConv := &automock.FetchRequestConverter{}
	defer mockFrConv.AssertExpectations(t)
	mockFrConv.On("InputFromGraphQL", mock.Anything).Return(nil)
	mockFrConv.On("ToGraphQL", mock.Anything).Return(nil)

	mockVersionConv := &automock.VersionConverter{}
	defer mockVersionConv.AssertExpectations(t)
	mockVersionConv.On("InputFromGraphQL", mock.Anything).Return(nil)
	mockVersionConv.On("ToGraphQL", mock.Anything).Return(nil)

	converter := api.NewConverter(mockAuthConv, mockFrConv, mockVersionConv)
	// WHEN & THEN
	convertedInputModel := converter.InputFromGraphQL(&graphql.APIDefinitionInput{Spec: &graphql.APISpecInput{}})
	require.NotNil(t, convertedInputModel)
	require.NotNil(t, convertedInputModel.Spec)
	require.Nil(t, convertedInputModel.Spec.Data)
	convertedAPIDef := convertedInputModel.ToAPIDefinition("id", "app_id")
	require.NotNil(t, convertedAPIDef)
	convertedGraphqlAPIDef := converter.ToGraphQL(convertedAPIDef)
	require.NotNil(t, convertedGraphqlAPIDef)
	assert.Nil(t, convertedGraphqlAPIDef.Spec.Data)
}

func TestEntityConverter_ToEntity(t *testing.T) {
	t.Run("success all nullable properites filled", func(t *testing.T) {
		//GIVEN
		apiModel := *fixDetailedModelAPIDefinition(uuid.New().String(), "name", "description", "group")
		versionConv := version.NewConverter()
		conv := api.NewConverter(nil, nil, versionConv)
		//WHEN
		entity, err := conv.ToEntity(apiModel)
		//THEN
		require.NoError(t, err)
		assertApiDefinition(t, apiModel, entity)
	})
	t.Run("success all nullable properites empty", func(t *testing.T) {
		//GIVEN
		apiModel := *fixModelAPIDefinition("id", "appid", "name", "desc")
		versionConv := version.NewConverter()
		conv := api.NewConverter(nil, nil, versionConv)
		//WHEN
		entity, err := conv.ToEntity(apiModel)
		//THEN
		require.NoError(t, err)
		assertApiDefinition(t, apiModel, entity)
	})
}

func TestEntityConverter_FromEntity(t *testing.T) {
	t.Run("success all nullable properties filled", func(t *testing.T) {
		//GIVEN
		entity := *fixDetailedApiDefinitionEntity("placeholder")
		versionConv := version.NewConverter()
		conv := api.NewConverter(nil, nil, versionConv)
		//WHEN
		apiModel, err := conv.FromEntity(entity)
		//THEN
		require.NoError(t, err)
		assertApiDefinition(t, apiModel, entity)
	})
	t.Run("success all nullable properties empty", func(t *testing.T) {
		//GIVEN
		entity := *fixMinimalApiDefinitionEntity("id", "app_id", "name", "target_url")
		versionConv := version.NewConverter()
		conv := api.NewConverter(nil, nil, versionConv)
		//WHEN
		apiModel, err := conv.FromEntity(entity)
		//THEN
		require.NoError(t, err)
		assertApiDefinition(t, apiModel, entity)
	})
}

func assertApiDefinition(t *testing.T, apiModel model.APIDefinition, entity api.APIDefinition) {
	assert.Equal(t, apiModel.ID, entity.ID)
	assert.Equal(t, apiModel.TenantID, entity.TenantID)
	assert.Equal(t, apiModel.ApplicationID, entity.AppID)
	assert.Equal(t, apiModel.Name, entity.Name)
	testdb.AssertSqlNullString(t, entity.Description, apiModel.Description)
	testdb.AssertSqlNullString(t, entity.Group, apiModel.Group)
	assert.Equal(t, apiModel.TargetURL, entity.TargetURL)
	assertAPISpec(t, entity.APISpec, apiModel.Spec)
	assertAuth(t, apiModel.DefaultAuth, entity.DefaultAuth)
	assertVersion(t, entity.Version, apiModel.Version)
}

func assertAPISpec(t *testing.T, specEntity *api.APISpec, apiSpec *model.APISpec) {
	if apiSpec != nil && specEntity != nil {
		testdb.AssertSqlNullString(t, specEntity.SpecData, apiSpec.Data)
		testdb.AssertSqlNullString(t, specEntity.SpecFormat, strings.Ptr(string(apiSpec.Format)))
		testdb.AssertSqlNullString(t, specEntity.SpecType, strings.Ptr(string(apiSpec.Type)))
	}
}

func assertAuth(t *testing.T, expected *model.Auth, defaultAuth sql.NullString) {
	var entityAuth *model.Auth

	if defaultAuth.Valid && defaultAuth.String != "" {
		entityAuth = &model.Auth{}
		err := json.Unmarshal([]byte(defaultAuth.String), entityAuth)
		require.NoError(t, err)
	}
	if expected != nil && entityAuth != nil {
		assert.Equal(t, expected.Credential, entityAuth.Credential)
		assert.Equal(t, expected.RequestAuth, entityAuth.RequestAuth)
		assert.Equal(t, expected.AdditionalQueryParams, entityAuth.AdditionalQueryParams)
		assert.Equal(t, expected.AdditionalHeaders, entityAuth.AdditionalHeaders)
	} else {
		assert.Nil(t, expected)
		assert.Nil(t, entityAuth)
	}
}

func assertVersion(t *testing.T, entity *version.Version, versionModel *model.Version) {
	if versionModel != nil {
		testdb.AssertSqlNullString(t, entity.VersionValue, &versionModel.Value)
		testdb.AssertSqlNullBool(t, entity.VersionDepracated, versionModel.Deprecated)
		testdb.AssertSqlNullString(t, entity.VersionDepracatedSince, versionModel.DeprecatedSince)
		testdb.AssertSqlNullBool(t, entity.VersionForRemoval, versionModel.ForRemoval)
	} else {
		assert.Nil(t, entity)
	}
}
