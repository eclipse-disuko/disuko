// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

// func TestShouldGetAllObligations(t *testing.T) {
// 	//given
// 	mockDataAccessLayer, handler := prepareObligationDataAccessLayerAndHandler(t)

// 	key := "key"
// 	mockDataAccessLayer.EXPECT().FindAll(requestSessionTest).Return([]obligation.Obligation{obligation.Obligation{Key: key, Name: "My Obligation"}}, nil).Times(1)

// 	//when
// 	recorder := RecordHTTPRequest(t, "GET", "/api/v1/admin/obligations", []byte(""), handler.GetAllHandler, nil)

// 	//then
// 	CheckHTTPCode(t, recorder, 200)

// 	var receivedResponse obligation.AllResponse
// 	json.Unmarshal(recorder.Body.Bytes(), &receivedResponse)

// 	assert.Len(t, receivedResponse.Obligation, 1, "expected number of item: %d", 1)
// 	assert.Equal(t, key, receivedResponse.Obligation[0].Key, "expected item key %s, but received %s", key, receivedResponse.Obligation[0].Key)
// }

// func prepareObligationDataAccessLayerAndHandler(t *testing.T) (*mocks.MockObligationAccessLayer, ObligationsHandler) {
// 	mockCtrl := gomock.NewController(t)
// 	mockDataAccessLayer := mocks.NewMockObligationAccessLayer(mockCtrl)
// 	handler := ObligationsHandler{ObligationRepository: mockDataAccessLayer}
// 	return mockDataAccessLayer, handler
// }
