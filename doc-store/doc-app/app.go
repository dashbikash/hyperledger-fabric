package main

import (
	"crypto/x509"
	"fmt"
	"os"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	mspID        = "Org1MSP"
	cryptoPath   = "../../fabric-samples/test-network/organizations/peerOrganizations/org1.example.com"
	certPath     = cryptoPath + "/users/User1@org1.example.com/msp/signcerts/cert.pem"
	keyPath      = cryptoPath + "/users/User1@org1.example.com/msp/keystore/bfc18bd08a9cf1918cdba9cef371fc922d0e41d7819c530320e70a0d9126a6aa_sk"
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

	// Submit transactions that store state to the ledger.
	// submitResult, err := contract.SubmitTransaction("CreateDocument", "arg1", "arg2")
	// panicOnError(err)
	// fmt.Printf("Submit result: %s", string(submitResult))

	// Evaluate transactions that query state from the ledger.
	evaluateResult, err := contract.EvaluateTransaction("GetAll")
	panicOnError(err)
	fmt.Printf("Evaluate result: %s", string(evaluateResult))
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
	certificatePEM, err := os.ReadFile(certPath)
	panicOnError(err)

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	panicOnError(err)

	id, err := identity.NewX509Identity(mspID, certificate)
	panicOnError(err)

	return id
}

// NewSign creates a function that generates a digital signature from a message digest using a private key.
func NewSign() identity.Sign {
	privateKeyPEM, err := os.ReadFile(keyPath)
	panicOnError(err)

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	panicOnError(err)

	sign, err := identity.NewPrivateKeySign(privateKey)
	panicOnError(err)

	return sign
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
