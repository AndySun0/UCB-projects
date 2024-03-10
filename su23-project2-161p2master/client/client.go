package client

// CS 161 Project 2

// Only the following imports are allowed! ANY additional imports
// may break the autograder!
// - bytes
// - encoding/hex
// - encoding/json
// - errors
// - fmt
// - github.com/cs161-staff/project2-userlib
// - github.com/google/uuid
// - strconv
// - strings

import (
	"bytes"

	"encoding/json"

	userlib "github.com/cs161-staff/project2-userlib"
	"github.com/google/uuid"

	// hex.EncodeToString(...) is useful for converting []byte to string

	// Useful for string manipulation
	"strings"

	// Useful for formatting strings (e.g. `fmt.Sprintf`).
	"fmt"

	// Useful for creating new error messages to return using errors.New("...")
	"errors"

	// Optional.
	_ "strconv"
)

// This serves two purposes: it shows you a few useful primitives,
// and suppresses warnings for imports not being used. It can be
// safely deleted!
func someUsefulThings() {

	// Creates a random UUID.
	randomUUID := uuid.New()

	// Prints the UUID as a string. %v prints the value in a default format.
	// See https://pkg.go.dev/fmt#hdr-Printing for all Golang format string flags.
	userlib.DebugMsg("Random UUID: %v", randomUUID.String())

	// Creates a UUID deterministically, from a sequence of bytes.
	hash := userlib.Hash([]byte("user-structs/alice"))
	deterministicUUID, err := uuid.FromBytes(hash[:16])
	if err != nil {
		// Normally, we would `return err` here. But, since this function doesn't return anything,
		// we can just panic to terminate execution. ALWAYS, ALWAYS, ALWAYS check for errors! Your
		// code should have hundreds of "if err != nil { return err }" statements by the end of this
		// project. You probably want to avoid using panic statements in your own code.
		panic(errors.New("An error occurred while generating a UUID: " + err.Error()))
	}
	userlib.DebugMsg("Deterministic UUID: %v", deterministicUUID.String())

	// Declares a Course struct type, creates an instance of it, and marshals it into JSON.
	type Course struct {
		name      string
		professor []byte
	}

	course := Course{"CS 161", []byte("Nicholas Weaver")}
	courseBytes, err := json.Marshal(course)
	if err != nil {
		panic(err)
	}

	userlib.DebugMsg("Struct: %v", course)
	userlib.DebugMsg("JSON Data: %v", courseBytes)

	// Generate a random private/public keypair.
	// The "_" indicates that we don't check for the error case here.
	var pk userlib.PKEEncKey
	var sk userlib.PKEDecKey
	pk, sk, _ = userlib.PKEKeyGen()
	userlib.DebugMsg("PKE Key Pair: (%v, %v)", pk, sk)

	// Here's an example of how to use HBKDF to generate a new key from an input key.
	// Tip: generate a new key everywhere you possibly can! It's easier to generate new keys on the fly
	// instead of trying to think about all of the ways a key reuse attack could be performed. It's also easier to
	// store one key and derive multiple keys from that one key, rather than
	originalKey := userlib.RandomBytes(16)
	derivedKey, err := userlib.HashKDF(originalKey, []byte("mac-key"))
	if err != nil {
		panic(err)
	}
	userlib.DebugMsg("Original Key: %v", originalKey)
	userlib.DebugMsg("Derived Key: %v", derivedKey)

	// A couple of tips on converting between string and []byte:
	// To convert from string to []byte, use []byte("some-string-here")
	// To convert from []byte to string for debugging, use fmt.Sprintf("hello world: %s", some_byte_arr).
	// To convert from []byte to string for use in a hashmap, use hex.EncodeToString(some_byte_arr).
	// When frequently converting between []byte and string, just marshal and unmarshal the data.
	//
	// Read more: https://go.dev/blog/strings

	// Here's an example of string interpolation!
	_ = fmt.Sprintf("%s_%d", "file", 1)
}

// This is the type definition for the User struct.
// A Go struct is like a Python or Java class - it can have attributes
// (e.g. like the Username attribute) and methods (e.g. like the StoreFile method below).
type User struct {
	Username     string               // This is the username of the user.
	passwordHash []byte               // This is the hash of the user's password.
	passwordSalt []byte               // This is the salt used to hash the user's password.
	FileKeys     map[string][]byte    // This is a map of fileame to encrypted file keys.
	FileUUIDs    map[string]uuid.UUID // This is a map of filename to file UUID in datastore
	Shares       map[string][]string  // This is a map of filename to list of usernames that have access to the file.
	Active       map[string]bool      // This is a map that associates each filename with a boolean indicating whether use has active access to the file.
	HmacKeys     map[string][]byte    // This is a map of filename to its HMAC keys.

	// You can add other attributes here if you want! But note that in order for attributes to
	// be included when this struct is serialized to/from JSON, they must be capitalized.
	// On the flipside, if you have an attribute that you want to be able to access from
	// this struct's methods, but you DON'T want that value to be included in the serialized value
	// of this struct that's stored in datastore, then you can use a "private" variable (e.g. one that
	// begins with a lowercase letter).
}
// type UserList struct {
// 	username []string         
// }
// This is the type definition for the AccessTreeNode and AccessTree struct.
type AccessTreeNode struct {
	Username string           // This is the username of the user.
	Children []AccessTreeNode // This is a list of the children node.
}

type AccessTree struct {
	Root AccessTreeNode // This is the root node of the access tree.
}

// This is the type definition for the File struct.
type File struct {
	Filename   string     // This is the filename of the file.
	Owner      string     // This is the username of the file's owner.
	Content    []byte     // This is the content of the file.
	AccessTree AccessTree // This is the access tree of the file.
}

// This is the type definition for the Invitation struct.
type Invitation struct {
	Sender string // This is the username of the invitation's creater.
	//FileUUID uuid.UUID // This is the UUID of the file that the invitation is for.
	//EncryptedFileKey []byte
	FileKey  []byte // This is the encrypted file key.
	Filename string // This is the filename of the file.
	HmacKey  []byte // This is the HMAC key of the file.
}

// This is the type definition for the FileKey struct.
//type FileKey struct {
//	FileUUID uuid.UUID
//	EncryptedFileKey []byte
//}

// Helper functions

func EncryptData(ek userlib.PKEEncKey, data []byte) (encData []byte, encKey []byte, iv []byte, err error) {
	// Generate random symmetric key and IV
	symKey := userlib.RandomBytes(16)
	iv = userlib.RandomBytes(16)

	// Encrypt data with symmetric key
	encData = userlib.SymEnc(symKey, iv, data)

	// Encrypt symmetric key with public key of recipient
	encKey, err = userlib.PKEEnc(ek, symKey)
	return
}

func DecryptData(dk userlib.PKEDecKey, encData []byte, encKey []byte, iv []byte) (data []byte, err error) {
	// Decrypt symmetric key with private key of recipient
	symKey, err := userlib.PKEDec(dk, encKey)
	if err != nil {
		return nil, err
	}

	// Decrypt data with symmetric key
	data = userlib.SymDec(symKey, encData)
	return
}

func GeneratePasswordHash(password string) (hash []byte, err error) {
	// Convert password to byte array
	pwdBytes := []byte(password)

	// Generate password hash
	hash = userlib.Argon2Key(pwdBytes, []byte("5f78a2c1b794e36d"), 16)///////////////////////////////////
	return
}

func GenerateEncryptionAndMACKeys(sourceKey []byte) (encKey []byte, macKey []byte, err error) {
	// Generate encryption key
	encKey, err = userlib.HashKDF(sourceKey, []byte("encryption"))
	if err != nil {
		return nil, nil, err
	}

	// Generate MAC key
	macKey, err = userlib.HashKDF(sourceKey, []byte("mac"))
	if err != nil {
		return nil, nil, err
	}

	return
}

// NOTE: The following methods have toy (insecure!) implementations.

func InitUser(username string, password string) (userdataptr *User, err error) {
	var userdata User
	userdata.Username = username
	// if userdata.passwordSalt == []byte("0000000000000000") {
	// 	userdata.passwordSalt = userlib.RandomBytes(16)}
	// print(userdata.passwordSalt, "000")
	
	userdata.passwordHash, err = GeneratePasswordHash(password)
	// print("paswod: ", userdata.passwordHash)
	// print(userdata.passwordSalt)// Initialize the maps
	userdata.FileKeys = make(map[string][]byte)
	userdata.FileUUIDs = make(map[string]uuid.UUID)
	userdata.Shares = make(map[string][]string)
	userdata.Active = make(map[string]bool)
	userdata.HmacKeys = make(map[string][]byte)
	hash := userlib.Hash([]byte(username))

	usernameUUID, err := uuid.FromBytes(hash[16:32])
	userlib.DatastoreSet(usernameUUID, userdata.passwordHash)

	// Marshal the user data and store it
	userDataBytes, err := json.Marshal(userdata)
	if err != nil {
		return nil, err
	}

	// Encrypt user data
	key, _, err := GenerateEncryptionAndMACKeys(userlib.Hash([]byte(password))[:16])

	if err != nil {
		return nil, err
	}
	key = key[0:16]
	// print(key) ///
	encryptedUserData := userlib.SymEnc(key,userlib.RandomBytes(16), userDataBytes ) 
	///

	userDataUUID := uuid.New()
	// Use the username to generate a deterministic UUID
	// hash := userlib.Hash([]byte(username))
	userDataUUID, err = uuid.FromBytes(hash[:16])
	if err != nil {
		return nil, err
	}
	userlib.DatastoreSet(userDataUUID, encryptedUserData)

	return &userdata, nil
}

func GetUser(username string, password string) (userdataptr *User, err error) {
	// var userdata User
	// Use the username to generate the UUID for the user data
	hash := userlib.Hash([]byte(username))
	
	userDataUUID, err := uuid.FromBytes(hash[:16])
	usernameUUID, err := uuid.FromBytes(hash[16:32])

	
	if err != nil {
		return nil, err
	}
	
	// Fetch the user data from the datastore
	encryptedUserDataBytes, ok := userlib.DatastoreGet(userDataUUID)
	
	if !ok {
		return nil, errors.New("user not found")
	}

	// Decrypt user data
	key, _, err := GenerateEncryptionAndMACKeys(userlib.Hash([]byte(password))[:16])
	if err != nil {
		return nil, err
	}
	userdatapasswordhash, ok := userlib.DatastoreGet(usernameUUID)
	key = key[0:16]

	userDataBytes := userlib.SymDec(key, encryptedUserDataBytes)

	// Unmarshal the user data
	var userdata User
	// tmppasswordSalt, ok := userlib.DatastoreGet(usernameUUID)
	err = json.Unmarshal(userDataBytes, &userdata)
	if err != nil {
		return nil, err
	}
	// Check the password
	userdata.passwordHash = userdatapasswordhash
	
	expectedPasswordHash, err := GeneratePasswordHash(password)

	if err != nil {
		return nil, err
	}
	if !bytes.Equal(userdata.passwordHash, expectedPasswordHash) {
		return nil, errors.New("invalid password")
	}


	return &userdata, nil
}

func (userdata *User) AppendToFile(filename string, content []byte) error {
	// Fetch the file content from the datastore
	fileUUID, ok := userdata.FileUUIDs[filename]
	if !ok {
		return errors.New(strings.ToTitle("file not found"))
	}
	encryptedContent, ok := userlib.DatastoreGet(fileUUID)
	
	if !ok {
		return errors.New(strings.ToTitle("file not found"))
	}

	// Decrypt the existing file content
	fileKey, ok := userdata.FileKeys[filename]

	if !ok {
		return errors.New(strings.ToTitle("file not found"))
	}
	iv := userlib.RandomBytes(16)
	
	existingContent := userlib.SymDec(fileKey, encryptedContent)
	

	// Append the new content to the existing content
	appendedContent := append(existingContent, content...)
	

	encryptedAppendedContent := userlib.SymEnc(fileKey, iv, appendedContent)

	// Store the encrypted appended content back to the datastore
	userlib.DatastoreSet(fileUUID, encryptedAppendedContent)
	// print("!!apppp!!!!!encryptedAppendedContent",":::",encryptedAppendedContent,"\n")
	// fmt.Println("!!apppp!!!!!encryptedAppendedContent",":::",encryptedAppendedContent,"\n")
	return nil
}

func (userdata *User) StoreFile(filename string, content []byte) (err error) {
	// Encrypt the file content
	fileKey := userlib.RandomBytes(16)
	iv := userlib.RandomBytes(16)

	// encryptedContent := append(tmpcontent,userlib.SymEnc(fileKey,iv,content)...)
	encryptedContent := userlib.SymEnc(fileKey,iv,content)//iv,
	

	// Store the encrypted file content
	fileUUID := uuid.New()
	userlib.DatastoreSet(fileUUID, encryptedContent)


	// Store the file metadata
	userdata.FileKeys[filename] = fileKey
	userdata.FileUUIDs[filename] = fileUUID
	userdata.HmacKeys[filename] = userlib.Hash(fileKey) // You will need to replace this with the correct HMAC key
	
	return nil
}

func (userdata *User) LoadFile(filename string) (content []byte, err error) {
	// Fetch the file content from the datastore
	fileUUID, ok := userdata.FileUUIDs[filename]
	if !ok {
		return nil, errors.New(strings.ToTitle("file not found"))
	}
	encryptedContent, ok := userlib.DatastoreGet(fileUUID)
	if !ok {
		return nil, errors.New(strings.ToTitle("file not found"))
	}
	// print("!!!!!!!!encryptedContent",":::", encryptedContent)
	// Decrypt the file content
	fileKey, ok := userdata.FileKeys[filename]

	if !ok {
		return nil, errors.New(strings.ToTitle("file not found"))
	}
	// fmt.Println("encryptedContent",":::",encryptedContent,"\n")
	content =userlib.SymDec(fileKey, encryptedContent)//[]16:
	// contentstring :=string(content)
	// fmt.Println("loadContentstring",":::",contentstring,"\n")
	// fmt.Println("loadContent",":::",content,"\n")

	return content, nil
}

func (userdata *User) CreateInvitation(filename string, recipientUsername string) (
	invitationPtr uuid.UUID, err error) {
	// fmt.Println("filekeyuser",":::",userdata.FileKeys,"\n")
	fileKey, ok := userdata.FileKeys[filename]
	// fmt.Println("filekey",":::",fileKey,"\n")
	if !ok {
		return uuid.Nil, errors.New("file not found!")
	}

	invitation := Invitation{
		Sender:   userdata.Username,
		FileKey:  fileKey,
		Filename: filename,
	}

	invitationBytes, err := json.Marshal(invitation)
	if err != nil {
		return uuid.Nil, err
	}

	invitationPtr = uuid.New()
	userlib.DatastoreSet(invitationPtr, invitationBytes)

	return invitationPtr, nil
}

func (userdata *User) AcceptInvitation(senderUsername string, invitationPtr uuid.UUID, filename string) error {
	invitationBytes, ok := userlib.DatastoreGet(invitationPtr)
	if !ok {
		return errors.New("invitation not found")
	}

	var invitation Invitation
	err := json.Unmarshal(invitationBytes, &invitation)
	if err != nil {
		return err
	}

	if invitation.Sender != senderUsername {
		return errors.New("invitation sender does not match")
	}

	userdata.FileKeys[filename] = invitation.FileKey
	userdata.FileUUIDs[filename] = invitationPtr

	return nil
}

func (userdata *User) RevokeAccess(filename string, recipientUsername string) error {
	fileUUID, ok := userdata.FileUUIDs[filename]
	if !ok {
		return errors.New("file not found")
	}

	newFileKey := userlib.RandomBytes(16)
	userdata.FileKeys[filename] = newFileKey

	content, err := userdata.LoadFile(filename)
	if err != nil {
		return err
	}

	iv := userlib.RandomBytes(userlib.AESBlockSizeBytes)
	ciphertext := userlib.SymEnc(newFileKey, iv, content)

	newFileUUID := uuid.New()
	userdata.FileUUIDs[filename] = newFileUUID

	userlib.DatastoreSet(newFileUUID, ciphertext)

	// Delete the old data from the datastore
	userlib.DatastoreDelete(fileUUID)

	return nil
}

func getAccessTreeRoot(filename string, ownerUsername string) (root AccessTreeNode, err error) {
	// This function retrieves the root of the access tree for a file.
	// It does so by hashing the filename and owner's username to get a UUID, then using this UUID to get the root from the datastore.
	hash := userlib.Hash([]byte(filename + ownerUsername))
	uuid, err := uuid.FromBytes(hash[:16])
	if err != nil {
		return AccessTreeNode{}, err
	}
	data, ok := userlib.DatastoreGet(uuid)
	if !ok {
		return AccessTreeNode{}, errors.New("file not found")
	}
	err = json.Unmarshal(data, &root)
	if err != nil {
		return AccessTreeNode{}, err
	}
	return root, nil
}

func updateAccessTreeRoot(filename string, ownerUsername string, root AccessTreeNode) (err error) {
	// This function updates the root of the access tree for a file.
	// It does so by hashing the filename and owner's username to get a UUID, then using this UUID to store the new root in the datastore.
	hash := userlib.Hash([]byte(filename + ownerUsername))
	uuid, err := uuid.FromBytes(hash[:16])
	if err != nil {
		return err
	}
	data, err := json.Marshal(root)
	if err != nil {
		return err
	}
	userlib.DatastoreSet(uuid, data)
	return nil
}

func removeNodeFromTree(root AccessTreeNode, username string) (newRoot AccessTreeNode, found bool) {
	// This function removes a node from the access tree.
	// If the node is the root, it returns the first child as the new root.
	// If the node is not the root, it recursively searches the tree and removes the node from its parent's children.
	// It returns a boolean indicating whether the node was found and removed.
	if root.Username == username {
		return root.Children[0], true
	}
	for i, child := range root.Children {
		newChild, found := removeNodeFromTree(child, username)
		if found {
			root.Children[i] = newChild
			return root, true
		}
	}
	return root, false
}
