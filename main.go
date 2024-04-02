package main

import (
	"context"
	"fmt"
	"github.com/ghodss/yaml"
	grpcRuntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	rolemanager "grpc-gnostic-openapi-v3-example/gen"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type Service struct {
	rolemanager.UnimplementedRoleManagerServer
}

func (s *Service) CreateRole(ctx context.Context, in *rolemanager.CreateRole_Request) (*rolemanager.CreateRole_Response, error) {
	return &rolemanager.CreateRole_Response{SampleBodyField: "hello world!"}, nil
}

func main() {
	server := grpc.NewServer()
	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		err = server.Serve(listener)
		if err != nil {
			log.Fatal(err)
		}
	}()

	rolemanager.RegisterRoleManagerServer(server, &Service{})

	conn, err := grpc.DialContext(context.Background(), "localhost:9090", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	gatewayMux := grpcRuntime.NewServeMux()
	err = rolemanager.RegisterRoleManagerHandler(context.Background(), gatewayMux, conn)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/openapi", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fromYaml("gen/openapi.yaml")))
	})
	mux.Handle("/", gatewayMux)

	errCh := make(chan error, 1)
	go func() {
		log.Println("server is running on port 8080")
		errCh <- http.ListenAndServe(":8080", mux)
	}()

	go func() {
		url := "http://localhost:8080/openapi"
		if runtime.GOOS == "darwin" {
			cmd := exec.Command("open", url)
			cmd.Start()
		} else {
			// 다른 운영 체제에서는 표준 출력으로 URL을 출력하고 사용자가 수동으로 열도록 유도할 수 있음
			println("Please open your browser and go to:", url)
		}
	}()
	<-errCh
	log.Fatal("server stopped")
}

// fromYaml reads a yaml file from srcPath and returns a string of html code that displays the swagger ui
func fromYaml(srcPath string) string {
	// check if srcPath exists
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		log.Fatal(err)
	}

	// check if src is a json file
	if filepath.Ext(srcPath) != ".yaml" {
		log.Fatal("srcPath is not a yaml file")
	}

	// read file from srcPath
	file, err := os.ReadFile(srcPath)
	if err != nil {
		log.Fatal(err)
	}

	json, err := yaml.YAMLToJSON(file)
	if err != nil {
		log.Fatal(err)
	}

	code := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@latest/swagger-ui.css">
    <style>
        html {
            box-sizing: border-box;
            overflow: -moz-scrollbars-vertical;
            overflow-y: scroll;
        }
        *,
        *:before,
        *:after {
            box-sizing: inherit;
        }
        body {
            margin: 0;
            background: #fafafa;
        }
    </style>
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@latest/swagger-ui-bundle.js"></script>
<script src="https://unpkg.com/swagger-ui-dist@latest/swagger-ui-standalone-preset.js"></script>
<script>
    window.onload = function() {
        window.ui = SwaggerUIBundle({
	        spec: %s,
            dom_id: '#swagger-ui',
            presets: [
                SwaggerUIBundle.presets.apis,
                SwaggerUIStandalonePreset
            ],
            layout: "StandaloneLayout",
            deepLinking: true
        });
    };
</script>
</body>
</html>
`, string(json))

	return code
}
