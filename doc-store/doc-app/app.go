package main

import (
	"crypto/x509"
	"net/http"
	"os"
	"path"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	mspID        = "Org1MSP"
	cryptoPath   = "../../fabric-samples/test-network/organizations/peerOrganizations/org1.example.com"
	certPath     = cryptoPath + "/users/User1@org1.example.com/msp/signcerts/"
	keyPath      = cryptoPath + "/users/User1@org1.example.com/msp/keystore/"
	tlsCertPath  = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
	peerEndpoint = "dns:///localhost:7051"
	gatewayPeer  = "peer0.org1.example.com"
)

func main() {
	// Create gRPC client connection, which should be shared by all gateway connections to this endpoint.
	clientConnection, err := NewGrpcConnection()
	panicOnError(err)
	defer clientConnection.Close()

	// Create client identity and signing implementation based on X.509 certificate and private key.
	id := NewIdentity()
	sign := NewSign()

	// Create a Gateway connection for a specific client identity.
	gateway, err := client.Connect(id, client.WithSign(sign), client.WithClientConnection(clientConnection))
	panicOnError(err)
	defer gateway.Close()

	// Obtain smart contract deployed on the network.
	network := gateway.GetNetwork("mychannel")
	contract := network.GetContract("docstore")

	router := httprouter.New()
	router.GET("/doc", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		resp := GetAll(contract)
		w.Write(resp)
	})
	router.GET("/doc/:id", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		resp := GetDocument(contract, p.ByName("id"))
		w.Write(resp)
	})
	router.POST("/doc", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		queryParams := r.URL.Query()
		resp := CreateDocument(contract, queryParams.Get("id"), r.URL.Query().Get("content"), "bikash")
		w.Write(resp)
	})
	router.DELETE("/doc/:id", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		resp := DeleteDocument(contract, p.ByName("id"))
		w.Write(resp)
	})

	http.ListenAndServe(":8080", router)

}

// NewGrpcConnection creates a new gRPC client connection
func NewGrpcConnection() (*grpc.ClientConn, error) {
	tlsCertificatePEM, err := os.ReadFile(tlsCertPath)
	panicOnError(err)

	tlsCertificate, err := identity.CertificateFromPEM(tlsCertificatePEM)
	panicOnError(err)

	certPool := x509.NewCertPool()
	certPool.AddCert(tlsCertificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, "")

	return grpc.NewClient(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
}

// NewIdentity creates a client identity for this Gateway connection using an X.509 certificate.
func NewIdentity() *identity.X509Identity {
	certificatePEM, err := readFirstFile(certPath)
	panicOnError(err)

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	panicOnError(err)

	id, err := identity.NewX509Identity(mspID, certificate)
	panicOnError(err)

	return id
}

// NewSign creates a function that generates a digital signature from a message digest using a private key.
func NewSign() identity.Sign {
	privateKeyPEM, err := readFirstFile(keyPath)
	panicOnError(err)

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	panicOnError(err)

	sign, err := identity.NewPrivateKeySign(privateKey)
	panicOnError(err)

	return sign
}

// Helper methods
func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
func readFirstFile(dirPath string) ([]byte, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}

	fileNames, err := dir.Readdirnames(1)
	if err != nil {
		return nil, err
	}

	return os.ReadFile(path.Join(dirPath, fileNames[0]))
}
