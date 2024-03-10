package client_test

// You MUST NOT change these default imports.  ANY additional imports may
// break the autograder and everyone will be sad.

import (
	// Some imports use an underscore to prevent the compiler from complaining
	// about unused imports.
	_ "encoding/hex"
	_ "errors"
	_ "strconv"
	_ "strings"
	"testing"

	// A "dot" import is used here so that the functions in the ginko and gomega
	// modules can be used without an identifier. For example, Describe() and
	// Expect() instead of ginko.Describe() and gomega.Expect().
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	userlib "github.com/cs161-staff/project2-userlib"

	"github.com/cs161-staff/project2-starter-code/client"
)

func TestSetupAndExecution(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Tests")
}

// ================================================
// Global Variables (feel free to add more!)
// ================================================
const defaultPassword = "password"
const emptyString = ""
const contentOne = "Bitcoin is Nick's favorite "
const contentTwo = "digital "
const contentThree = "cryptocurrency!"

// ================================================
// Describe(...) blocks help you organize your tests
// into functional categories. They can be nested into
// a tree-like structure.
// ================================================

var _ = Describe("Client Tests", func() {

	// A few user declarations that may be used for testing. Remember to initialize these before you
	// attempt to use them!
	var alice *client.User
	var bob *client.User
	var charles *client.User
	var doris *client.User
	// var eve *client.User
	// var frank *client.User
	// var grace *client.User
	// var horace *client.User
	// var ira *client.User

	// These declarations may be useful for multi-session testing.
	var alicePhone *client.User
	var aliceLaptop *client.User
	var aliceDesktop *client.User

	var err error

	// A bunch of filenames that may be useful.
	aliceFile := "aliceFile.txt"
	bobFile := "bobFile.txt"
	charlesFile := "charlesFile.txt"
	dorisFile := "dorisFile.txt"
	// eveFile := "eveFile.txt"
	// frankFile := "frankFile.txt"
	// graceFile := "graceFile.txt"
	// horaceFile := "horaceFile.txt"
	// iraFile := "iraFile.txt"

	BeforeEach(func() {
		// This runs before each test within this Describe block (including nested tests).
		// Here, we reset the state of Datastore and Keystore so that tests do not interfere with each other.
		// We also initialize
		userlib.DatastoreClear()
		userlib.KeystoreClear()
	})

	Describe("Basic Tests", func() {

		Specify("Basic Test: Testing InitUser/GetUser on a single user.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting user Alice.")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())
		})

		Specify("Basic Test: Testing Single User Store/Load/Append.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentTwo)
			err = alice.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentThree)
			err = alice.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Loading file...")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))
		})

		Specify("Basic Test: Testing Create/Accept Invite Functionality with multiple users and multiple instances.", func() {
			userlib.DebugMsg("Initializing users Alice (aliceDesktop) and Bob.")
			aliceDesktop, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting second instance of Alice - aliceLaptop")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop storing file %s with content: %s", aliceFile, contentOne)
			err = aliceDesktop.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceLaptop creating invite for Bob.")
			invite, err := aliceLaptop.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepting invite from Alice under filename %s.", bobFile)
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob appending to file %s, content: %s", bobFile, contentTwo)
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop appending to file %s, content: %s", aliceFile, contentThree)
			err = aliceDesktop.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that aliceDesktop sees expected file data.")
			data, err := aliceDesktop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that aliceLaptop sees expected file data.")
			data, err = aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that Bob sees expected file data.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Getting third instance of Alice - alicePhone.")
			alicePhone, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that alicePhone sees Alice's changes.")
			data, err = alicePhone.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))
		})

		Specify("Basic Test: Testing Revoke Functionality", func() {
			userlib.DebugMsg("Initializing users Alice, Bob, and Charlie.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice storing file %s with content: %s", aliceFile, contentOne)
			alice.StoreFile(aliceFile, []byte(contentOne))

			userlib.DebugMsg("Alice creating invite for Bob for file %s, and Bob accepting invite under name %s.", aliceFile, bobFile)

			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob can load the file.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Bob creating invite for Charles for file %s, and Charlie accepting invite under name %s.", bobFile, charlesFile)
			invite, err = bob.CreateInvitation(bobFile, "charles")
			Expect(err).To(BeNil())

			err = charles.AcceptInvitation("bob", invite, charlesFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Bob can load the file.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Charles can load the file.")
			data, err = charles.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Alice revoking Bob's access from %s.", aliceFile)
			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err = alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob/Charles lost access to the file.")
			_, err = bob.LoadFile(bobFile)
			Expect(err).ToNot(BeNil())

			_, err = charles.LoadFile(charlesFile)
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Checking that the revoked users cannot append to the file.")
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())

			err = charles.AppendToFile(charlesFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())
		})

	})
	Describe("My Tests Cases", func() {
		//#1
		Specify("Case#1", func() {
			userlib.DebugMsg("Initialising user alice")
			_, err := client.InitUser("alice", "password123")
			Expect(err).ShouldNot(HaveOccurred())
			userlib.DebugMsg("Trying to initialize alice again")
			_, err = client.InitUser("alice", "password321")
			Expect(err).Should(HaveOccurred())
		})
		//#2
		Specify("Case#2", func() {
			userlib.DebugMsg("Starting user Alice")
			alice, err := client.InitUser("alice", defaultPassword)
			Expect(err).ShouldNot(HaveOccurred())
			userlib.DebugMsg("Saving file content: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).ShouldNot(HaveOccurred())
			userlib.DebugMsg("Retrieving file content")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data).Should(Equal([]byte(contentOne)))
		})
		//#3
		Specify("Case#3", func() {
			userlib.DebugMsg("Starting user Bob without password")
			_, err := client.InitUser("Bob", "")
			Expect(err).ShouldNot(HaveOccurred())
		})
		//#4
		Specify("Case#4", func() {
			userlib.DebugMsg("Initiating user Bob")
			_, err := client.InitUser("Bob", "password123")
			Expect(err).ShouldNot(HaveOccurred())
			userlib.DebugMsg("Initiating another user bob")
			_, err = client.InitUser("bob", "password123")
			Expect(err).ShouldNot(HaveOccurred())
		})
		//#5
		Specify("Case#5", func() {
			userlib.DebugMsg("Initiating user Bob")
			bob, err := client.InitUser("Bob", "bestpassword321")
			Expect(err).ShouldNot(HaveOccurred())
			userlib.DebugMsg("Initiating user Eve")
			eve, err := client.InitUser("Eve", "bestpassword321")
			Expect(err).ShouldNot(HaveOccurred())

			eve.StoreFile("file123", []byte("abc"))
			inv, err := eve.CreateInvitation("file123", "Bob")
			Expect(err).ShouldNot(HaveOccurred())

			err = bob.AcceptInvitation("Eve", inv, "file321")
			Expect(err).ShouldNot(HaveOccurred())

			err = bob.AppendToFile("file321", []byte("cba"))
			Expect(err).ShouldNot(HaveOccurred())

			byteContent, err := eve.LoadFile("file123")
			Expect(err).ShouldNot(HaveOccurred())
			content := string(byteContent)
			Expect(content).Should(Equal("abccba"))
		})
		//#6
		Specify("Case#6", func() {
			userlib.DebugMsg("Initiating users Bob, Alice, and Eve")
			bob, err := client.InitUser("Bob", "bestpassword321")
			Expect(err).ShouldNot(HaveOccurred())
			eve, err := client.InitUser("Eve", "bestpassword321")
			Expect(err).ShouldNot(HaveOccurred())
			alice, err := client.InitUser("Alice", "bestpassword321")
			Expect(err).ShouldNot(HaveOccurred())

			userlib.DebugMsg("Bob creates a document.")
			const FILENAME = "document"
			const CONTENT = "A neutral document about nothing in particular."
			bob.StoreFile(FILENAME, []byte(CONTENT))

			userlib.DebugMsg("Bob sends invite to Eve")
			inv, err := bob.CreateInvitation(FILENAME, "Eve")
			Expect(err).ShouldNot(HaveOccurred())

			userlib.DebugMsg("Eve acknowledges the invite")
			eve.AcceptInvitation("Bob", inv, FILENAME)

			userlib.DebugMsg("Eve sends invite to Alice")
			inv, err = eve.CreateInvitation(FILENAME, "Alice")
			Expect(err).ShouldNot(HaveOccurred())

			userlib.DebugMsg("Alice acknowledges the invite")
			alice.AcceptInvitation("Alice", inv, FILENAME)

			userlib.DebugMsg("Bob tries to revoke access for Alice")
			err = bob.RevokeAccess(FILENAME, "Alice")
			Expect(err).ToNot(BeNil())
		})
		//#7
		Specify("Case#7", func() {
			userlib.DebugMsg("Initiating user Bob")
			bob, err := client.InitUser("Bob", "bestpassword321")
			Expect(err).ShouldNot(HaveOccurred())

			userlib.DebugMsg("Creating file with no name")
			err = bob.StoreFile("", []byte(""))
			Expect(err).ShouldNot(HaveOccurred())

		})
		//#8
		Specify("Case#8", func() {
			userlib.DebugMsg("Initiating user Bob")
			bob, err := client.InitUser("Bob", "preferredpassword123")
			Expect(err).ShouldNot(HaveOccurred())

			const TEXT = `
        	Testing tex: During the whimsical adventure in the land of make-believe, the blibberjam and the snickerdoodle engaged in a lively game of quibblesnack, accompanied by the melodious tunes of the zoodlepluff orchestra. As the jibberflap danced gracefully under the moonlight, a group of babbledorfs gathered to share tales of their latest wobblefuzz escapades. Meanwhile, a ziggledorf and a snugglewump indulged in a picnic on rainbow-hued meadows, enjoying slices of ziggledorf pie and giggling uncontrollably. Amidst the laughter and merriment, the wobblefuzz and the quibblesnack orchestrated a hopscotch tournament on fluffy marshmallow clouds, causing even the most serious of jibberflaps to burst into fits of laughter. The enchanted day concluded with a harmonious choir of snickerdoodles singing enchanting melodies that echoed through the zoodlepluff forest, leaving memories of nonsensical joy in their wake.
        `
			const DOC_NAME = "documentName"

			userlib.DatastoreResetBandwidth()
			userlib.DebugMsg("Bob saves a substantial file")
			err = bob.StoreFile(DOC_NAME, []byte(TEXT))
			Expect(err).ShouldNot(HaveOccurred())

			bandwidthInitial := userlib.DatastoreGetBandwidth()
			userlib.DatastoreResetBandwidth()

			userlib.DebugMsg("Appending to file data")
			err = bob.AppendToFile(DOC_NAME, []byte("."))
			Expect(err).ShouldNot(HaveOccurred())

			bandwidthPostAppend := userlib.DatastoreGetBandwidth()

			Expect(bandwidthInitial > bandwidthPostAppend).To(BeTrue())
		})
		//#9
		Specify("Case#9", func() {
			userlib.DebugMsg("Initiating users Bob and Alice")
			bob, err := client.InitUser("Bob", "preferredpassword123")
			Expect(err).ShouldNot(HaveOccurred())
			alice, err := client.InitUser("Alice", "preferredpassword123")
			Expect(err).ShouldNot(HaveOccurred())

			const DOC_NAME = "documentName"
			const TEXT_0 = "Quite fascinating"
			const TEXT_1 = "Rather fascinating"
			const TEXT_2 = "Extremely fascinating"

			userlib.DebugMsg("Bob saves a file")
			err = bob.StoreFile(DOC_NAME, []byte(TEXT_0))
			Expect(err).ShouldNot(HaveOccurred())

			userlib.DebugMsg("Bob shares file with Alice")
			inv, err := bob.CreateInvitation(DOC_NAME, "Alice")
			Expect(err).ShouldNot(HaveOccurred())
			alice.AcceptInvitation("Bob", inv, DOC_NAME)

			userlib.DebugMsg("Alice retrieves file")
			byteContent, err := alice.LoadFile(DOC_NAME)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(byteContent)).To(Equal(TEXT_0))

			userlib.DebugMsg("Bob re-saves the file with the same name")
			err = bob.StoreFile(DOC_NAME, []byte(TEXT_1))
			Expect(err).ShouldNot(HaveOccurred())

			userlib.DebugMsg("Alice retrieves the file again")
			byteContent, err = alice.LoadFile(DOC_NAME)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(byteContent)).To(Equal(TEXT_1))

			userlib.DebugMsg("Alice saves a different file under the same name")
			err = alice.StoreFile(DOC_NAME, []byte(TEXT_2))
			Expect(err).ShouldNot(HaveOccurred())

			userlib.DebugMsg("Bob retrieves the file yet again")
			byteContent, err = bob.LoadFile(DOC_NAME)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(byteContent)).To(Equal(TEXT_2))
		})

		// //#10
		Specify("Case#10", func() {
			userlib.DebugMsg("Initiating users Bob and Alice")
			bob, err := client.InitUser("Bob", "preferredpassword123")
			Expect(err).ShouldNot(HaveOccurred())
			alice, err := client.InitUser("Alice", "preferredpassword123")
			Expect(err).ShouldNot(HaveOccurred())

			const DOC_NAME = "documentName"
			const TEXT_1 = "Remarkable"
			const TEXT_2 = "Ordinary"

			userlib.DebugMsg("Bob saves a file with a particular name")
			err = bob.StoreFile(DOC_NAME, []byte(TEXT_1))
			Expect(err).ShouldNot(HaveOccurred())

			userlib.DebugMsg("Alice saves another file with an identical name")
			err = alice.StoreFile(DOC_NAME, []byte(TEXT_2))
			Expect(err).ShouldNot(HaveOccurred())

			userlib.DebugMsg("Alice retrieves her specific file")
			byteContent, err := alice.LoadFile(DOC_NAME)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(byteContent)).To(Equal(TEXT_2))

		})

		// #11
		Specify("Case#11 ;", func() {
			userlib.DebugMsg("Setting up Alice (aliceMachine), Bob, Charlie, and Doris.")
			aliceMachine, e := client.InitUser("alice", defaultPassword)
			Expect(e).To(BeNil())

			bobUser, e := client.InitUser("bob", defaultPassword)
			Expect(e).To(BeNil())

			charlieUser, e := client.InitUser("charles", defaultPassword)
			Expect(e).To(BeNil())

			dorisUser, e := client.InitUser("doris", defaultPassword)
			Expect(e).To(BeNil())

			userlib.DebugMsg("Getting second instance of Alice - aliceMobile")
			aliceMobile, e := client.GetUser("alice", defaultPassword)
			Expect(e).To(BeNil())

			userlib.DebugMsg("aliceMachine storing: %s with %s", aliceFile, contentOne)
			e = aliceMachine.StoreFile(aliceFile, []byte(contentOne))
			Expect(e).To(BeNil())

			userlib.DebugMsg("aliceMobile sharing with Bob.")
			sharedInv, e := aliceMobile.CreateInvitation(aliceFile, "bob")
			Expect(e).To(BeNil())

			userlib.DebugMsg("Bob accepts invite for file as %s.", bobFile)
			e = bobUser.AcceptInvitation("alice", sharedInv, bobFile)
			Expect(e).To(BeNil())

			userlib.DebugMsg("aliceMobile shares with Doris.")
			sharedInvAD, e := aliceMobile.CreateInvitation(aliceFile, "doris")
			Expect(e).To(BeNil())

			userlib.DebugMsg("Doris accepts invite for file as %s.", dorisFile)
			e = dorisUser.AcceptInvitation("alice", sharedInvAD, dorisFile)
			Expect(e).To(BeNil())

			userlib.DebugMsg("Bob shares with Charles.")
			sharedInvBC, e := bobUser.CreateInvitation(bobFile, "charles")
			Expect(e).To(BeNil())

			userlib.DebugMsg("Charles accepts invite for file as %s.", charlesFile)
			e = charlieUser.AcceptInvitation("bob", sharedInvBC, charlesFile)
			Expect(e).To(BeNil())

			userlib.DebugMsg("Bob updates file %s with: %s", bobFile, contentTwo)
			e = bobUser.AppendToFile(bobFile, []byte(contentTwo))
			Expect(e).To(BeNil())

			userlib.DebugMsg("aliceMachine appends to file %s with: %s", aliceFile, contentThree)
			e = aliceMachine.AppendToFile(aliceFile, []byte(contentThree))
			Expect(e).To(BeNil())

			userlib.DebugMsg("Checking that aliceDesktop sees expected file data.")
			data, err := aliceDesktop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that Bob sees expected file data.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that Charles sees expected file data.")
			data, err = charles.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that Doris sees expected file data.")
			data, err = doris.LoadFile(dorisFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))
		})

		// #12
		Specify("Basic Test: Testing Revoke Functionality", func() {
			userlib.DebugMsg("Setting up Alice, Bob, Charlie, and Doris for revocation tests.")
			aliceUser, e := client.InitUser("alice", defaultPassword)
			Expect(e).To(BeNil())

			bobUser, e := client.InitUser("bob", defaultPassword)
			Expect(e).To(BeNil())

			charlieUser, e := client.InitUser("charles", defaultPassword)
			Expect(e).To(BeNil())

			dorisUser, e := client.InitUser("doris", defaultPassword)
			Expect(e).To(BeNil())

			userlib.DebugMsg("Alice saves file %s: %s", aliceFile, contentOne)
			aliceUser.StoreFile(aliceFile, []byte(contentOne))

			invForBob, e := aliceUser.CreateInvitation(aliceFile, "bob")
			userlib.DebugMsg("Alice shares with Bob for %s, and Bob accepts as: %s.", aliceFile, bobFile)
			Expect(e).To(BeNil())
			e = bobUser.AcceptInvitation("alice", invForBob, bobFile)
			Expect(e).To(BeNil())

			userlib.DebugMsg("Alice creating invite for Doris for file %s, and Doris accepting invite under name %s.", aliceFile, dorisFile)

			inviteAD, err := alice.CreateInvitation(aliceFile, "doris")
			Expect(err).To(BeNil())

			err = doris.AcceptInvitation("alice", inviteAD, dorisFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob can load the file.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Bob creating invite for Charles for file %s, and Charlie accepting invite under name %s.", bobFile, charlesFile)
			invite, err = bob.CreateInvitation(bobFile, "charles")
			Expect(err).To(BeNil())

			err = charles.AcceptInvitation("bob", invite, charlesFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Bob can load the file.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Charles can load the file.")
			data, err = charles.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Alice revoking Bob's access from %s.", aliceFile)
			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err = alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob/Charles lost access to the file.")
			_, err = bob.LoadFile(bobFile)
			Expect(err).ToNot(BeNil())

			_, err = charles.LoadFile(charlesFile)
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Checking that Doris can still load the file.")
			data, err = doris.LoadFile(dorisFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that the revoked users cannot append to the file.")
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())

			err = charles.AppendToFile(charlesFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())
		})
		//#13
		Specify("Case13", func() {
			userlib.DebugMsg("Initializing Alice, Bob, Charlie, and Doris.")
			aliceUser, e := client.InitUser("alice", defaultPassword)
			Expect(e).To(BeNil())

			bobUser, e := client.InitUser("bob", defaultPassword)
			Expect(e).To(BeNil())

			charlieUser, e := client.InitUser("charles", defaultPassword)
			Expect(e).To(BeNil())

			dorisUser, e := client.InitUser("doris", defaultPassword)
			Expect(e).To(BeNil())

			userlib.DebugMsg("Alice saving file %s: %s", aliceFile, contentOne)
			aliceUser.StoreFile(aliceFile, []byte(contentOne))

			inv, e := aliceUser.CreateInvitation(aliceFile, "bob")
			userlib.DebugMsg("Alice's invitation for Bob for file: %s, Bob accepts as: %s.", aliceFile, bobFile)
			Expect(e).To(BeNil())
			e = bobUser.AcceptInvitation("alice", inv, bobFile)
			Expect(e).To(BeNil())

			invForDoris, e := aliceUser.CreateInvitation(aliceFile, "doris")
			userlib.DebugMsg("Revoking Doris's permission from %s before she accepts.", aliceFile)
			Expect(e).To(BeNil())
			e = aliceUser.RevokeAccess(aliceFile, "doris")
			Expect(e).To(BeNil())

			e = dorisUser.AcceptInvitation("alice", invForDoris, dorisFile)
			Expect(e).ToNot(BeNil())

			userlib.DebugMsg("Verifying Doris's file access.")
			_, e = dorisUser.LoadFile(dorisFile)
			Expect(e).ToNot(BeNil())

			userlib.DebugMsg("Validating Alice's file access.")
			content, e := aliceUser.LoadFile(aliceFile)
			Expect(e).To(BeNil())
			Expect(content).To(Equal([]byte(contentOne)))
		})
	})
})
