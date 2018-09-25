package main

import  (
    "crypto/sha256"
    "fmt"
    "pbc"
)

type messageData struct {
    message   string
    signature []byte
}

func main() {
    // System parameters
    params := pbc.GenerateA(160, 512)
    pairing := params.NewPairing()
    g := pairing.NewG2().Rand()

    // Params and g to Alice and Bob
    sharedParams := params.String()
    sharedG := g.Bytes()

    // Channel for messages (maybe network connection)
    messageChannel := make(chan *messageData)

    // Channel for public key distribution. 
    keyChannel := make(chan []byte)

    // Channel to wait until both simulations are done
    finished := make(chan bool)

    // Simulate the conversation participants
    go alice(sharedParams, sharedG, messageChannel, keyChannel, finished)
    go bob(sharedParams, sharedG, messageChannel, keyChannel, finished)

    // Wait for the communication to finish
    <-finished
    <-finished

}

// Alice 
func alice(sharedParams string, sharedG []byte, messageChannel chan *messageData, keyChannel chan []byte, finished chan bool) {
    // Alice loads the system parameters
    pairing, _ := pbc.NewPairingFromString(sharedParams)
    g := pairing.NewG2().SetBytes(sharedG)

    // Generate keypair (x, g^x)
    privKey := pairing.NewZr().Rand()
    fmt.Println("Alice PrivKey")
    fmt.Println(privKey)
    pubKey := pairing.NewG2().PowZn(g, privKey)
    fmt.Println("Alice Pubkey")
    fmt.Println(pubKey)
    // Send public key to Bob
    keyChannel <- pubKey.Bytes()

    // Some time later, sign a message, hashed to h, as h^x
    message := "Hello, Fantom Foundation!"
    h := pairing.NewG1().SetFromStringHash(message, sha256.New())
    fmt.Println("Message: " + message)
    fmt.Println("Alice signature")
    signature := pairing.NewG2().PowZn(h, privKey)
    fmt.Println(signature)
    
    // Send the message and signature to Bob
    messageChannel <- &messageData{message: message, signature: signature.Bytes()}

    finished <- true
}

// Bob 
func bob(sharedParams string, sharedG []byte, messageChannel chan *messageData, keyChannel chan []byte, finished chan bool) {
    // Load system parameters
    pairing, _ := pbc.NewPairingFromString(sharedParams)
    g := pairing.NewG2().SetBytes(sharedG)

    // Bob receives Alice's public key 
    pubKey := pairing.NewG2().SetBytes(<-keyChannel)
    fmt.Println("Bob Pubkey")
    fmt.Println(pubKey)

    // Bob receives a message to verify
    data := <-messageChannel
    signature := pairing.NewG1().SetBytes(data.signature)
    fmt.Println("Bob signature")
    fmt.Println(signature)
    
    // To verify, check that e(h,g^x)=e(sig,g)
    h := pairing.NewG1().SetFromStringHash(data.message, sha256.New())
    temp1 := pairing.NewGT().Pair(h, pubKey)
    fmt.Println("e(h, g^x)")
    fmt.Println(temp1)
    temp2 := pairing.NewGT().Pair(signature, g)
    fmt.Println("e(sig, g)")
    fmt.Println(temp2)
    if !temp1.Equals(temp2) {
        fmt.Println("*BUG* Signature check failed *BUG*")
    } else {
        fmt.Println("Signature verified correctly - e(h,g^x)=e(sig,g)")
    }

    finished <- true
}
