package jwt

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	j "github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
)

func TestOpts_Unmarshall(t *testing.T) {
	type args struct {
		desc    string
		opts    func() IOts
		rewrite func(token string) string
		wantJti string
		wantErr string
	}

	uuid := "b9094e37-7018-4412-9137-525d739dac65"
	signaturErr := "signature is invalid"

	tests := []args{
		{
			desc: "when ok",
			opts: func() IOts {
				o := MustBuildOpts(t)
				o.Jti = uuid
				return o
			},
			rewrite: func(token string) string {
				// nothing replaced
				return MustReplacePayload(t, token, jti, uuid)
			},
			wantJti: uuid,
		}, {
			desc: "when algo differs",
			opts: func() IOts {
				o := MustBuildOpts(t)
				o.Algo = HS256
				return o
			},
			rewrite: func(token string) string {
				return MustReplaceHeader(t, token, "alg", HS512.Alg())
			},
			wantErr: ErrAlgoMismatch.Error(),
		}, {
			desc: "when jti differs",
			opts: func() IOts {
				o := MustBuildOpts(t)
				o.Jti = uuid
				return o
			},
			rewrite: func(token string) string {
				return MustReplacePayload(t, token, jti, "foo")
			},
			wantErr: signaturErr,
		}, {
			desc: "when aud differs",
			opts: func() IOts {
				o := MustBuildOpts(t)
				o.Aud = "aud"
				return o
			},
			rewrite: func(token string) string {
				return MustReplacePayload(t, token, aud, "foo")
			},
			wantErr: signaturErr,
		}, {
			desc: "when sub differs",
			opts: func() IOts {
				o := MustBuildOpts(t)
				o.Sub = "sub"
				return o
			},
			rewrite: func(token string) string {
				return MustReplacePayload(t, token, sub, "foo")
			},
			wantErr: signaturErr,
		}, {
			desc: "when iss differs",
			opts: func() IOts {
				o := MustBuildOpts(t)
				o.Sub = "iss"
				return o
			},
			rewrite: func(token string) string {
				return MustReplacePayload(t, token, iss, "foo")
			},
			wantErr: signaturErr,
		}, {
			desc: "when exp differs",
			opts: func() IOts {
				o := MustBuildOpts(t)
				return o
			},
			rewrite: func(token string) string {
				return MustReplacePayload(t, token, exp, 1_0000_000_000)
			},
			wantErr: signaturErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			token, err := tt.opts().Marshall()
			require.NoError(t, err)

			gor, err := tt.opts().Unmarshall(tt.rewrite(token))
			if tt.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantJti, gor)
			}
		})
	}
}

func MustBuildOpts(t *testing.T) Opts {
	return Opts{
		Jti:  "jti",
		Aud:  "aud",
		Sub:  "sub",
		Iss:  "iss",
		Algo: HS256,
		Key:  []byte{1, 2, 3},
		Exp:  time.Now().Add(1 * time.Second),
	}
}

func MustReplacePayload(t *testing.T, token string, key string, value interface{}) string {
	return MustReplace(t, token, "", "", key, value)
}

func MustReplaceHeader(t *testing.T, token string, key string, value interface{}) string {
	return MustReplace(t, token, key, value, "", "")
}

func MustReplace(t *testing.T, token string, hk string, hv interface{}, pk string, pv interface{}) string {
	var headerMap map[string]interface{}
	var payloadMap map[string]interface{}

	sep := "."
	parts := strings.Split(token, sep)
	headerB64, payloadB64, sigB64 := parts[0], parts[1], parts[2]

	header, err := j.DecodeSegment(headerB64)
	require.NoError(t, err)

	payload, err := j.DecodeSegment(payloadB64)
	require.NoError(t, err)

	err = json.Unmarshal(header, &headerMap)
	require.NoError(t, err)

	err = json.Unmarshal(payload, &payloadMap)
	require.NoError(t, err)

	if hk != "" {
		headerMap[hk] = hv
	}

	if pk != "" {
		payloadMap[pk] = pv
	}

	header, err = json.Marshal(headerMap)
	require.NoError(t, err)

	payload, err = json.Marshal(payloadMap)
	require.NoError(t, err)

	return strings.Join([]string{
		j.EncodeSegment(header),
		j.EncodeSegment(payload),
		sigB64,
	}, sep)
}
