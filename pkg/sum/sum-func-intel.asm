
00000000004f9aa0 <github.com/efficientgo/examples/pkg/sum.Sum>:
  4f9aa0:	49 3b 66 10          	cmp    rsp,QWORD PTR [r14+0x10]
  4f9aa4:	0f 86 e6 00 00 00    	jbe    4f9b90 <github.com/efficientgo/examples/pkg/sum.Sum+0xf0>
  4f9aaa:	48 83 ec 70          	sub    rsp,0x70
  4f9aae:	48 89 6c 24 68       	mov    QWORD PTR [rsp+0x68],rbp
  4f9ab3:	48 8d 6c 24 68       	lea    rbp,[rsp+0x68]
  4f9ab8:	48 89 44 24 78       	mov    QWORD PTR [rsp+0x78],rax
  4f9abd:	90                   	nop
  4f9abe:	66 90                	xchg   ax,ax
  4f9ac0:	e8 db 64 f8 ff       	call   47ffa0 <os.ReadFile>
  4f9ac5:	48 85 ff             	test   rdi,rdi
  4f9ac8:	75 35                	jne    4f9aff <github.com/efficientgo/examples/pkg/sum.Sum+0x5f>
  4f9aca:	c6 44 24 47 0a       	mov    BYTE PTR [rsp+0x47],0xa
  4f9acf:	48 8d 7c 24 47       	lea    rdi,[rsp+0x47]
  4f9ad4:	be 01 00 00 00       	mov    esi,0x1
  4f9ad9:	49 89 f0             	mov    r8,rsi
  4f9adc:	45 31 c9             	xor    r9d,r9d
  4f9adf:	49 c7 c2 ff ff ff ff 	mov    r10,0xffffffffffffffff
  4f9ae6:	e8 35 52 fb ff       	call   4aed20 <bytes.genSplit>
  4f9aeb:	48 85 db             	test   rbx,rbx
  4f9aee:	7e 0b                	jle    4f9afb <github.com/efficientgo/examples/pkg/sum.Sum+0x5b>
  4f9af0:	48 89 5c 24 48       	mov    QWORD PTR [rsp+0x48],rbx
  4f9af5:	31 c9                	xor    ecx,ecx
  4f9af7:	31 d2                	xor    edx,edx
  4f9af9:	eb 33                	jmp    4f9b2e <github.com/efficientgo/examples/pkg/sum.Sum+0x8e>
  4f9afb:	31 c0                	xor    eax,eax
  4f9afd:	eb 12                	jmp    4f9b11 <github.com/efficientgo/examples/pkg/sum.Sum+0x71>
  4f9aff:	31 c0                	xor    eax,eax
  4f9b01:	48 89 fb             	mov    rbx,rdi
  4f9b04:	48 89 f1             	mov    rcx,rsi
  4f9b07:	48 8b 6c 24 68       	mov    rbp,QWORD PTR [rsp+0x68]
  4f9b0c:	48 83 c4 70          	add    rsp,0x70
  4f9b10:	c3                   	ret    
  4f9b11:	31 db                	xor    ebx,ebx
  4f9b13:	31 c9                	xor    ecx,ecx
  4f9b15:	48 8b 6c 24 68       	mov    rbp,QWORD PTR [rsp+0x68]
  4f9b1a:	48 83 c4 70          	add    rsp,0x70
  4f9b1e:	c3                   	ret    
  4f9b1f:	4c 8b 44 24 60       	mov    r8,QWORD PTR [rsp+0x60]
  4f9b24:	49 8d 40 18          	lea    rax,[r8+0x18]
  4f9b28:	48 89 d1             	mov    rcx,rdx
  4f9b2b:	48 89 f2             	mov    rdx,rsi
  4f9b2e:	48 89 4c 24 58       	mov    QWORD PTR [rsp+0x58],rcx
  4f9b33:	48 89 44 24 60       	mov    QWORD PTR [rsp+0x60],rax
  4f9b38:	48 89 54 24 50       	mov    QWORD PTR [rsp+0x50],rdx
  4f9b3d:	48 8b 18             	mov    rbx,QWORD PTR [rax]
  4f9b40:	48 8b 70 08          	mov    rsi,QWORD PTR [rax+0x8]
  4f9b44:	31 c0                	xor    eax,eax
  4f9b46:	48 89 f1             	mov    rcx,rsi
  4f9b49:	e8 52 43 f5 ff       	call   44dea0 <runtime.slicebytetostring>
  4f9b4e:	b9 0a 00 00 00       	mov    ecx,0xa
  4f9b53:	bf 40 00 00 00       	mov    edi,0x40
  4f9b58:	e8 43 8d f8 ff       	call   4828a0 <strconv.ParseInt>
  4f9b5d:	0f 1f 00             	nop    DWORD PTR [rax]
  4f9b60:	48 85 db             	test   rbx,rbx
  4f9b63:	75 1f                	jne    4f9b84 <github.com/efficientgo/examples/pkg/sum.Sum+0xe4>
  4f9b65:	48 8b 54 24 58       	mov    rdx,QWORD PTR [rsp+0x58]
  4f9b6a:	48 ff c2             	inc    rdx
  4f9b6d:	48 8b 74 24 50       	mov    rsi,QWORD PTR [rsp+0x50]
  4f9b72:	48 01 c6             	add    rsi,rax
  4f9b75:	48 8b 7c 24 48       	mov    rdi,QWORD PTR [rsp+0x48]
  4f9b7a:	48 39 d7             	cmp    rdi,rdx
  4f9b7d:	7f a0                	jg     4f9b1f <github.com/efficientgo/examples/pkg/sum.Sum+0x7f>
  4f9b7f:	48 89 f0             	mov    rax,rsi
  4f9b82:	eb 8d                	jmp    4f9b11 <github.com/efficientgo/examples/pkg/sum.Sum+0x71>
  4f9b84:	31 c0                	xor    eax,eax
  4f9b86:	48 8b 6c 24 68       	mov    rbp,QWORD PTR [rsp+0x68]
  4f9b8b:	48 83 c4 70          	add    rsp,0x70
  4f9b8f:	c3                   	ret    
  4f9b90:	48 89 44 24 08       	mov    QWORD PTR [rsp+0x8],rax
  4f9b95:	48 89 5c 24 10       	mov    QWORD PTR [rsp+0x10],rbx
  4f9b9a:	e8 61 85 f6 ff       	call   462100 <runtime.morestack_noctxt.abi0>
  4f9b9f:	48 8b 44 24 08       	mov    rax,QWORD PTR [rsp+0x8]
  4f9ba4:	48 8b 5c 24 10       	mov    rbx,QWORD PTR [rsp+0x10]
  4f9ba9:	e9 f2 fe ff ff       	jmp    4f9aa0 <github.com/efficientgo/examples/pkg/sum.Sum>

