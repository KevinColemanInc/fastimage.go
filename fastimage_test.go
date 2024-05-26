package fastimage

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		httpClient HttpClient
	}
	tests := []struct {
		name string
		args args
		want *fastimageImpl
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.httpClient); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fastimageImpl_ValidateURL(t *testing.T) {
	type fields struct {
		httpStatus    string
		contentLength int64
	}
	type args struct {
		url    string
		policy *Policy
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantOk       bool
		wantBytes    []byte
		wantResponse *http.Response
		wantErr      bool
	}{
		{
			name: "smoke",
			fields: fields{
				httpStatus:    http.StatusText(http.StatusOK),
				contentLength: 500,
			},
			args: args{
				url:    "https://www.kcoleman.me/images/sf-bridge.jpg",
				policy: NewPolicy().WithImageTypes(JPEG),
			},
			wantOk:    true,
			wantBytes: []byte{255, 216, 255, 224, 0, 16, 74, 70, 73, 70, 0, 1, 1},
			wantErr:   false,
		},
		{
			name: "smoke - webp",
			fields: fields{
				httpStatus:    http.StatusText(http.StatusOK),
				contentLength: 500,
			},
			args: args{
				url:    "https://lh3.googleusercontent.com/a/ACg8ocIVNjFOrBp3ZJJHgdwYwKdily-Y9OCGx1BO0McD0SQSX97eUD2d=s83-c-mo",
				policy: NewPolicy().WithImageTypes(PNG, JPEG, WEBP),
			},
			wantOk:    true,
			wantBytes: []byte{255, 216, 255, 224, 0, 16, 74, 70, 73, 70, 0, 1, 1},
			wantErr:   false,
		},
		{
			name: "invalid - html as image",
			args: args{
				url:    "http://example.com/image.jpg",
				policy: NewPolicy().WithImageTypes(JPEG),
			},
			wantOk:    false,
			wantBytes: []byte{},
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//body := io.NopCloser(bytes.NewReader(tt.wantBytes))
			//mockHttpClient := HttpClientMock{
			//	DoFunc: func(req *http.Request) (*http.Response, error) {
			//		return &http.Response{
			//			Status:        tt.fields.httpStatus,
			//			ContentLength: tt.fields.contentLength,
			//			Body:          body,
			//		}, nil
			//	},
			//}
			i := &fastimageImpl{
				httpClient: http.DefaultClient,
			}
			got, got1, _, err := i.ValidateURL(tt.args.url, tt.args.policy)
			if (err != nil) != tt.wantErr {
				t.Errorf("fastimageImpl.ValidateURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantOk {
				t.Errorf("fastimageImpl.ValidateURL() got = %v, want %v", got, tt.wantOk)
			}
			if !reflect.DeepEqual(got1, tt.wantBytes) {
				t.Errorf("fastimageImpl.ValidateURL() got1 = %v, want %v", got1, tt.wantBytes)
			}
		})
	}
}
