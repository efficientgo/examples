TEXT github.com/efficientgo/examples/pkg/sum.Sum(SB) /home/bwplotka/Repos/examples/pkg/sum/sum.go
func Sum(fileName string) (ret int64, _ error) {
  0x4f9aa0		493b6610		CMPQ 0x10(R14), SP	
  0x4f9aa4		0f86e6000000		JBE 0x4f9b90		
  0x4f9aaa		4883ec70		SUBQ $0x70, SP		
  0x4f9aae		48896c2468		MOVQ BP, 0x68(SP)	
  0x4f9ab3		488d6c2468		LEAQ 0x68(SP), BP	
  0x4f9ab8		4889442478		MOVQ AX, 0x78(SP)	
	b, err := ioutil.ReadFile(fileName)
  0x4f9abd		90			NOPL			
  0x4f9abe		6690			NOPW			
	return os.ReadFile(filename)
  0x4f9ac0		e8db64f8ff		CALL os.ReadFile(SB)	
	if err != nil {
  0x4f9ac5		4885ff			TESTQ DI, DI		
  0x4f9ac8		7535			JNE 0x4f9aff		
	for _, line := range bytes.Split(b, []byte("\n")) {
  0x4f9aca		c64424470a		MOVB $0xa, 0x47(SP)	
func Split(s, sep []byte) [][]byte { return genSplit(s, sep, 0, -1) }
  0x4f9acf		488d7c2447		LEAQ 0x47(SP), DI	
  0x4f9ad4		be01000000		MOVL $0x1, SI		
  0x4f9ad9		4989f0			MOVQ SI, R8		
  0x4f9adc		4531c9			XORL R9, R9		
  0x4f9adf		49c7c2ffffffff		MOVQ $-0x1, R10		
  0x4f9ae6		e83552fbff		CALL bytes.genSplit(SB)	
	for _, line := range bytes.Split(b, []byte("\n")) {
  0x4f9aeb		4885db			TESTQ BX, BX		
  0x4f9aee		7e0b			JLE 0x4f9afb		
  0x4f9af0		48895c2448		MOVQ BX, 0x48(SP)	
  0x4f9af5		31c9			XORL CX, CX		
  0x4f9af7		31d2			XORL DX, DX		
  0x4f9af9		eb33			JMP 0x4f9b2e		
  0x4f9afb		31c0			XORL AX, AX		
  0x4f9afd		eb12			JMP 0x4f9b11		
		return 0, err
  0x4f9aff		31c0			XORL AX, AX		
  0x4f9b01		4889fb			MOVQ DI, BX		
  0x4f9b04		4889f1			MOVQ SI, CX		
  0x4f9b07		488b6c2468		MOVQ 0x68(SP), BP	
  0x4f9b0c		4883c470		ADDQ $0x70, SP		
  0x4f9b10		c3			RET			
	return ret, nil
  0x4f9b11		31db			XORL BX, BX		
  0x4f9b13		31c9			XORL CX, CX		
  0x4f9b15		488b6c2468		MOVQ 0x68(SP), BP	
  0x4f9b1a		4883c470		ADDQ $0x70, SP		
  0x4f9b1e		c3			RET			
	for _, line := range bytes.Split(b, []byte("\n")) {
  0x4f9b1f		4c8b442460		MOVQ 0x60(SP), R8	
  0x4f9b24		498d4018		LEAQ 0x18(R8), AX	
  0x4f9b28		4889d1			MOVQ DX, CX		
		ret += num
  0x4f9b2b		4889f2			MOVQ SI, DX		
	for _, line := range bytes.Split(b, []byte("\n")) {
  0x4f9b2e		48894c2458		MOVQ CX, 0x58(SP)	
  0x4f9b33		4889442460		MOVQ AX, 0x60(SP)	
		ret += num
  0x4f9b38		4889542450		MOVQ DX, 0x50(SP)	
	for _, line := range bytes.Split(b, []byte("\n")) {
  0x4f9b3d		488b18			MOVQ 0(AX), BX		
  0x4f9b40		488b7008		MOVQ 0x8(AX), SI	
		num, err := strconv.ParseInt(string(line), 10, 64)
  0x4f9b44		31c0			XORL AX, AX				
  0x4f9b46		4889f1			MOVQ SI, CX				
  0x4f9b49		e85243f5ff		CALL runtime.slicebytetostring(SB)	
  0x4f9b4e		b90a000000		MOVL $0xa, CX				
  0x4f9b53		bf40000000		MOVL $0x40, DI				
  0x4f9b58		e8438df8ff		CALL strconv.ParseInt(SB)		
  0x4f9b5d		0f1f00			NOPL 0(AX)				
		if err != nil {
  0x4f9b60		4885db			TESTQ BX, BX		
  0x4f9b63		751f			JNE 0x4f9b84		
	for _, line := range bytes.Split(b, []byte("\n")) {
  0x4f9b65		488b542458		MOVQ 0x58(SP), DX	
  0x4f9b6a		48ffc2			INCQ DX			
		ret += num
  0x4f9b6d		488b742450		MOVQ 0x50(SP), SI	
  0x4f9b72		4801c6			ADDQ AX, SI		
	for _, line := range bytes.Split(b, []byte("\n")) {
  0x4f9b75		488b7c2448		MOVQ 0x48(SP), DI	
  0x4f9b7a		4839d7			CMPQ DX, DI		
  0x4f9b7d		7fa0			JG 0x4f9b1f		
	return ret, nil
  0x4f9b7f		4889f0			MOVQ SI, AX		
  0x4f9b82		eb8d			JMP 0x4f9b11		
			return 0, err
  0x4f9b84		31c0			XORL AX, AX		
  0x4f9b86		488b6c2468		MOVQ 0x68(SP), BP	
  0x4f9b8b		4883c470		ADDQ $0x70, SP		
  0x4f9b8f		c3			RET			
func Sum(fileName string) (ret int64, _ error) {
  0x4f9b90		4889442408		MOVQ AX, 0x8(SP)					
  0x4f9b95		48895c2410		MOVQ BX, 0x10(SP)					
  0x4f9b9a		e86185f6ff		CALL runtime.morestack_noctxt.abi0(SB)			
  0x4f9b9f		488b442408		MOVQ 0x8(SP), AX					
  0x4f9ba4		488b5c2410		MOVQ 0x10(SP), BX					
  0x4f9ba9		e9f2feffff		JMP github.com/efficientgo/examples/pkg/sum.Sum(SB)	
