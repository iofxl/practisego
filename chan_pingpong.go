package main

import "time"
import "fmt"

func main() {
	var Ball int
	table := make(chan int)

	for i := 0; i < 100; i++ {
		go player(i, table)
	}

	table <- Ball

	time.Sleep(1 * time.Second)
	fmt.Println(<-table)
}

func player(who int, table chan int) {
	for {
		ball := <-table
		ball++
		time.Sleep(100 * time.Millisecond)
		fmt.Println(who, ball)
		table <- ball
	}
}

/* result:

1 1
0 2
6 3
3 4
4 5
5 6
8 7
7 8
9 9
10 10
2 11
19 12
11 13
15 14
13 15
14 16
17 17
16 18
12 19
23 20
18 21
21 22
22 23
26 24
20 25
24 26
28 27
27 28
25 29
29 30
30 31
32 32
33 33
31 34
35 35
34 36
37 37
38 38
39 39
36 40
49 41
40 42
41 43
42 44
43 45
44 46
45 47
46 48
47 49
48 50
59 51
50 52
55 53
54 54
51 55
52 56
56 57
53 58
57 59
66 60
58 61
60 62
63 63
61 64
62 65
71 66
64 67
65 68
67 69
68 70
73 71
72 72
70 73
74 74
69 75
75 76
88 77
76 78
77 79
78 80
79 81
80 82
81 83
82 84
83 85
84 86
85 87
86 88
87 89
91 90
89 91
90 92
92 93
97 94
95 95
96 96
99 97
98 98
94 99
93 100
1 101
0 102
6 103
3 104
4 105
5 106
8 107
7 108
9 109
109
*/
