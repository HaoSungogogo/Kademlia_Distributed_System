For this project we are allowing groups to earn up to 10 points of extra credit by creating their own unit tests or by improving the thoroughness of the provided tests.

In order to receive extra credit, groups *must* include a plain text file called "EXTRACREDIT" in the "libkademlia" directory. Students must also include a clear explanation of what each new or improved function is testing. The list of functions to be considered for E.C. must go in the "EXTRACREDIT" file. Explanations can either go in the "EXTRACREDIT" file or in the comments of the group's *_test.go* file. You can either add to the provided test file, or create a new one (just be clear in your explanation).

Each new "test" of functionality can earn between 0 and 5 points, depending on the complexity and usefulness of the test. Note that a individual function can "test" multiple aspects of Kademlia, as the provided test functions do. However, each aspect of functionality that is being tested must be explained.

Tests that are, for the most part, redundant with the existing tests will not receive credit. Groups can earn a maximum of 10 points total.



We tested following two additional aspects regarding to assignment 1:

1. store conflict:

  In this case, we test the system's ability to detect hash-key conflict. We tried to store two key-value pair with same key value, The system is expected to detect the key conflict and return an error.

  Please see the test function on custom_test.go/TestStoreConflict

  Our system worked as expected to report the key conflict error, However, on testing this aspect: We found a incomplete interface on golang's rpc package. The rpc can not pass error format string, and will report [rpc: gob error encoding body: gob: type not registered for interface: errors.errorString] error on passing back error string on client.Call() command.

  We confirm this is rpc's implementation problem on this [source]( https://groups.google.com/forum/#!topic/golang-dev/Cua1Av1J8Nc).

  We comment out the rpc bugging function, line 78 - 83 in file custom_test.go. You can restore the commented code to see the system could return error, but the rpc can not handle the error properly.

2. retrieving more than K nodes:

  In this case, we test the system's ability to find more than K node on target contact.

  Please see the test function on custom_test.go/TestFindKNodes

  This is a rather easy test, we linked Node B to 30 contacts (including Node A), and then, from Node A, we send findNode command to Node B.

  The system is expected to return 20 contacts, as findNode should return K contacts when it's possible.
