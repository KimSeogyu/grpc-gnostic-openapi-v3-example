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

	postgwpb "github.com/KimSeogyu/grpc-gnostic-openapi-v3-example/internal/proto/gateway/postpb/v1"
	rolemanagergwpb "github.com/KimSeogyu/grpc-gnostic-openapi-v3-example/internal/proto/gateway/rolemanagerpb/v1"
	"github.com/KimSeogyu/grpc-gnostic-openapi-v3-example/internal/proto/postpb/v1"
	"github.com/KimSeogyu/grpc-gnostic-openapi-v3-example/internal/proto/rolemanagerpb/v1"
)

type Service struct {
	rolemanagerpbv1.UnimplementedV1RoleManagerServiceServer
	postpbv1.UnimplementedV1PostServiceServer
}

func (s *Service) CreateRole(ctx context.Context, in *rolemanagerpbv1.CreateRoleRequest) (*rolemanagerpbv1.CreateRoleResponse, error) {
	return &rolemanagerpbv1.CreateRoleResponse{SampleBodyField: "hello world!"}, nil
}

func (s *Service) CreatePost(ctx context.Context, in *postpbv1.CreatePostRequest) (*postpbv1.CreatePostResponse, error) {
	return &postpbv1.CreatePostResponse{SampleBodyField: "hello world!"}, nil
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

	rolemanagerpbv1.RegisterV1RoleManagerServiceServer(server, &Service{})
	postpbv1.RegisterV1PostServiceServer(server, &Service{})

	conn, err := grpc.DialContext(context.Background(), "localhost:9090", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	gatewayMux := grpcRuntime.NewServeMux()
	err = rolemanagergwpb.RegisterV1RoleManagerServiceHandler(context.Background(), gatewayMux, conn)
	if err != nil {
		log.Fatal(err)
	}

	err = postgwpb.RegisterV1PostServiceHandler(context.Background(), gatewayMux, conn)
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
