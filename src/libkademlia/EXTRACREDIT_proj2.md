We implemented the required function together with three test function, listed as follow:

1. TestIterativeFindNode()
2. TestIterativeStore()
3. TestIterativeFindValue()

All three test case is implemented in proj2_test.go

All of the following three testing function used a tree structure similar like the tree in TestFindNode.

TestIterativeFindNode:

In this test we try to DoIterativeFindNode from one of the leaf node tree_node[0], try to targeting another leaf node tree_node[3] using it's id.

Notice that the leaf nodes do not know each other in the first place, but if the iterative find node function work properly, they should eventually find and connect to each other. And should return a list of all possible nearby node's id, ideally with the tree_node[3]'s id on top, as it is the closest to target id.

On the test, our function can return a list as expected, with tree_node[3]'s id on top.

TestIterativeStore:

Our Testing function test to store value "hello world" with tree_node[3]'s id as key value, from tree_node[0], ideally this key-value pair will be stored to all possible nearby node all together.

After we called the IterativeStore, we try to locally find value from each leaf node, if the function worked properly, all the leaf node should contain the key-value pair.

Our test show that our IterativeStore can work as expected.

TestIterativeFindValue:

In this test, we start with call a DoStoreValue from the root of the tree using tree_node[3]'s id as key, so the key-value pair should be only stored in tree_node[3].

Then we initiate IterativeFindValue from tree_node[0], notice tree_node[0] do not know tree_node[3] in this stage, if our function worked properly, tree_node[0] should perform a similar action as IterativeFindNode, until the key-value pair is founded in tree_node[3].

In our test, our function can work as expected and eventually find the key-value pair in tree_node[3].
