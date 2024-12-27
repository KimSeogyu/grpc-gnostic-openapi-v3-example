package main

import (
	"context"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"

	openapi "github.com/KimSeogyu/grpc-gnostic-openapi-v3-example/docs"
	grpcRuntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	postgwpb "github.com/KimSeogyu/grpc-gnostic-openapi-v3-example/internal/proto/gateway/postpb"
	rolemanagergwpb "github.com/KimSeogyu/grpc-gnostic-openapi-v3-example/internal/proto/gateway/rolemanagerpb"
	"github.com/KimSeogyu/grpc-gnostic-openapi-v3-example/internal/proto/postpb"
	"github.com/KimSeogyu/grpc-gnostic-openapi-v3-example/internal/proto/rolemanagerpb"
)

type Service struct {
	rolemanagerpb.UnimplementedRoleManagerServer
	postpb.UnimplementedPostServer
}

func (s *Service) CreateRole(ctx context.Context, in *rolemanagerpb.CreateRole_Request) (*rolemanagerpb.CreateRole_Response, error) {
	return &rolemanagerpb.CreateRole_Response{SampleBodyField: "hello world!"}, nil
}

func (s *Service) CreatePost(ctx context.Context, in *postpb.CreatePost_Request) (*postpb.CreatePost_Response, error) {
	return &postpb.CreatePost_Response{SampleBodyField: "hello world!"}, nil
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

	rolemanagerpb.RegisterRoleManagerServer(server, &Service{})
	postpb.RegisterPostServer(server, &Service{})

	conn, err := grpc.DialContext(context.Background(), "localhost:9090", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	gatewayMux := grpcRuntime.NewServeMux()
	err = rolemanagergwpb.RegisterRoleManagerHandler(context.Background(), gatewayMux, conn)
	if err != nil {
		log.Fatal(err)
	}

	err = postgwpb.RegisterPostHandler(context.Background(), gatewayMux, conn)
	if err != nil {
		log.Fatal(err)
	}

	serveFileFS := func(w http.ResponseWriter, r *http.Request, fsys fs.FS, name string) {
		fs := http.FileServer(http.FS(fsys))
		r.URL.Path = name
		fs.ServeHTTP(w, r)
	}
	err = gatewayMux.HandlePath("GET", "/spec", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		serveFileFS(w, r, openapi.StaticFiles, "spec.html")
	})
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", gatewayMux)

	errCh := make(chan error, 1)
	go func() {
		log.Println("server is running on port 8080")
		errCh <- http.ListenAndServe(":8080", mux)
	}()

	go func() {
		url := "http://localhost:8080/spec"
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
