package node

import (
	"errors"
	"io"
	"testing"

	"github.com/docker/cli/internal/test"
	"github.com/docker/cli/internal/test/builders"
	"github.com/docker/docker/api/types/swarm"
	"gotest.tools/v3/assert"
)

func TestNodePromoteErrors(t *testing.T) {
	testCases := []struct {
		args            []string
		nodeInspectFunc func() (swarm.Node, []byte, error)
		nodeUpdateFunc  func(nodeID string, version swarm.Version, node swarm.NodeSpec) error
		expectedError   string
	}{
		{
			expectedError: "requires at least 1 argument",
		},
		{
			args: []string{"nodeID"},
			nodeInspectFunc: func() (swarm.Node, []byte, error) {
				return swarm.Node{}, []byte{}, errors.New("error inspecting the node")
			},
			expectedError: "error inspecting the node",
		},
		{
			args: []string{"nodeID"},
			nodeUpdateFunc: func(nodeID string, version swarm.Version, node swarm.NodeSpec) error {
				return errors.New("error updating the node")
			},
			expectedError: "error updating the node",
		},
	}
	for _, tc := range testCases {
		cmd := newPromoteCommand(
			test.NewFakeCli(&fakeClient{
				nodeInspectFunc: tc.nodeInspectFunc,
				nodeUpdateFunc:  tc.nodeUpdateFunc,
			}))
		cmd.SetArgs(tc.args)
		cmd.SetOut(io.Discard)
		cmd.SetErr(io.Discard)
		assert.ErrorContains(t, cmd.Execute(), tc.expectedError)
	}
}

func TestNodePromoteNoChange(t *testing.T) {
	cmd := newPromoteCommand(
		test.NewFakeCli(&fakeClient{
			nodeInspectFunc: func() (swarm.Node, []byte, error) {
				return *builders.Node(builders.Manager()), []byte{}, nil
			},
			nodeUpdateFunc: func(nodeID string, version swarm.Version, node swarm.NodeSpec) error {
				if node.Role != swarm.NodeRoleManager {
					return errors.New("expected role manager, got" + string(node.Role))
				}
				return nil
			},
		}))
	cmd.SetArgs([]string{"nodeID"})
	assert.NilError(t, cmd.Execute())
}

func TestNodePromoteMultipleNode(t *testing.T) {
	cmd := newPromoteCommand(
		test.NewFakeCli(&fakeClient{
			nodeInspectFunc: func() (swarm.Node, []byte, error) {
				return *builders.Node(), []byte{}, nil
			},
			nodeUpdateFunc: func(nodeID string, version swarm.Version, node swarm.NodeSpec) error {
				if node.Role != swarm.NodeRoleManager {
					return errors.New("expected role manager, got" + string(node.Role))
				}
				return nil
			},
		}))
	cmd.SetArgs([]string{"nodeID1", "nodeID2"})
	assert.NilError(t, cmd.Execute())
}
