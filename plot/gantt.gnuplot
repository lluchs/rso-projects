#set terminal pngcairo size 75cm,75cm
set terminal svg size 2000,2000

set lmargin 40

set xdata time
timeformat = "%Y-%m-%d"
set format x "%b\n'%y"
set datafile separator "\t"

set yrange [-1:]
OneMonth = strptime("%m","2")
set xtics OneMonth nomirror
set xtics scale 2, 0.5
set mxtics 4
#set ytics nomirror
unset ytics
set grid x y
unset key
set title "{/=15 RSO Projects}"
set border 1

T(N) = timecolumn(N,timeformat)

set style arrow 1 filled size screen 0.01, 15 fixed lt 3 lw 1.5

plot "gantt_all.dat" using (T(2)) : ($0) : (T(3)-T(2)) : (0.0) : yticlabel(1) with vector as 1, \
     "gantt_all.dat" using (T(2)) : ($0) : 1 with labels right offset -2

# vim: ft=gnuplot
