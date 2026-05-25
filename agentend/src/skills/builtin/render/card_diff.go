package main

import (
	"crypto/rand"
	"fmt"
)

// cmdDiff 输出 diff 卡片标记，携带 snapshotId（UUID v4）
func cmdDiff() {
	fmt.Printf("```%s\ntype: diff\nsnapshotId: %s\n```\n", blockMarker(), newUUIDv4())
}

func newUUIDv4() string {
	var buf [16]byte
	if _, err := rand.Read(buf[:]); err != nil {
		panic(err)
	}
	buf[6] = (buf[6] & 0x0f) | 0x40 // version 4
	buf[8] = (buf[8] & 0x3f) | 0x80 // variant 10
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		buf[0:4], buf[4:6], buf[6:8], buf[8:10], buf[10:16])
}
