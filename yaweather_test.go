package yaweather_test

//
// func fixture(t *testing.T) string {
// 	t.Helper()
// 	data, err := os.ReadFile("testdata/forecast.html")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	return string(data)
// }
//
// type transportFunc func(*http.Request) (*http.Response, error)
//
// func (f transportFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
//
// func TestFetchForecasts(t *testing.T) {
// 	previous := http.DefaultClient
// 	t.Cleanup(func() { http.DefaultClient = previous })
// 	location := yaweather.Location{Latitude: 55.75, Longitude: 37.62}
// 	page := fixture(t)
// 	for _, tc := range []struct {
// 		name   string
// 		status int
// 		body   string
// 		cancel bool
// 		fail   bool
// 	}{
// 		{"success", http.StatusOK, page, false, false},
// 		{"http error", http.StatusServiceUnavailable, page, false, true},
// 		{"invalid page", http.StatusOK, "<html>captcha</html>", false, true},
// 		{"oversized page", http.StatusOK, strings.Repeat("x", (8<<20)+1), false, true},
// 		{"canceled", http.StatusOK, page, true, true},
// 	} {
// 		t.Run(tc.name, func(t *testing.T) {
// 			http.DefaultClient = &http.Client{Transport: transportFunc(func(r *http.Request) (*http.Response, error) {
// 				if r.Method != http.MethodGet || r.URL.String() != location.URL() {
// 					t.Fatalf("wrong forecast request: %s %s", r.Method, r.URL)
// 				}
// 				if deadline, ok := r.Context().Deadline(); !ok || time.Until(deadline) > 10*time.Second {
// 					t.Fatal("missing forecast timeout")
// 				}
// 				if err := r.Context().Err(); err != nil {
// 					return nil, err
// 				}
// 				return &http.Response{StatusCode: tc.status, Body: io.NopCloser(strings.NewReader(tc.body))}, nil
// 			})}
// 			ctx, cancel := context.WithCancel(t.Context())
// 			defer cancel()
// 			if tc.cancel {
// 				cancel()
// 			}
// 			forecasts, err := location.GetForecastsContext(ctx)
// 			if (err != nil) != tc.fail {
// 				t.Fatalf("error = %v", err)
// 			}
// 			if tc.cancel && !errors.Is(err, context.Canceled) {
// 				t.Fatalf("lost cancellation: %v", err)
// 			}
// 			if !tc.fail && (len(forecasts) != 10 || forecasts[0].Day.Condition == yaweather.ConditionUnknown) {
// 				t.Fatalf("invalid forecasts: %+v", forecasts)
// 			}
// 		})
// 	}
// }
