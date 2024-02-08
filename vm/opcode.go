package vm

type OpCode = int

const (
	OP_LOAD_GLOBAL  = "LOAD_GLOBAL"
	OP_LOAD_CONST   = "LOAD_CONST"
	OP_POP_CONST    = "POP_CONST"
	OP_LOAD_VAR     = "LOAD_VAR"
	FUNC_CALL       = "FUNC_CALL"
	OP_PUSH_CONST   = "PUSH_CONST"
	OP_PUSH_SCOPE   = "PUSH_SCOPE"
	OP_POP_SCOPE    = "POP_SCOPE"
	OP_NOT          = "NOT"
	OP_NEGATIVE     = "OP_NEGATIVE"
	OP_STORE_GLOBAL = "STORE_GLOBAL"
	OP_DECL_VAR     = "DECL_VAR"
	OP_SET_VAR      = "SET_VAR"
	OP_MULT         = "MULT"
	OP_EQ           = "EQ"
	OP_CMP          = "CMP"
	OP_ADD          = "ADD"
	OP_DIV          = "DIV"
	OP_SUB          = "SUB"
	OP_RETURN       = "RETURN"
	OP_HALT         = "HALT"
	OP_DEBUG_PRINT  = "DEBUG_PRINT"
	OP_LABEL        = "LABEL"
	JUMP_NOT        = "JUMP_NOT"
	JUMP_IF         = "JUMP_IF"
	JUMP            = "JUMP"
)
