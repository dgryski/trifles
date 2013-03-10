#!/usr/bin/perl

# assume 80 col screen
$keycols = 8;
$starcols = 50;
$cntcols = 8;

# keep track of the values
while (<>) { chomp; push @values, $_; ++$hist{$_}; $n++}

# the most common value
$maxcnt = 0; $mode = 0;
foreach (keys %hist) {
   if ($maxcnt < $hist{$_}) {
      $mode = $_;
      $maxcnt = $hist{$_};
   }
}

# mean
$sum = 0;
for (@values) { $sum += $_; }
$avg = $sum / $n;

# standard deviation
$sdev = 0;
for (@values) { $sdev += ($_-$avg)**2 / $n; }
$sdev = sqrt($sdev);

# print historgram

# header
$fmt = "%${keycols}s  %-${starcols}s  %${cntcols}s\n";
printf ($fmt, "Key", "", "Count");

# separator
print "-" x ($keycols + $starcols + $cntcols + 4), "\n";

# heuristic -- sort numerically or not?
if (grep (!/\d/, keys %hist)) {
   # yup, some text strings
   @hkeys = sort { $a cmp $b } keys %hist;
} else {
   # nope, digits all the way
   @hkeys = sort { $a <=> $b } keys %hist;
}

# produce the graph
# and find the mode
foreach (@hkeys) {
    $ast = '*' x ($starcols * ($hist{$_} / $maxcnt));
#    $ast = '*' x ($starcols * $hist{$_} / $n);
    printf $fmt, $_, $ast, $hist{$_};
}

# find the median, which we can do now 'cause the keys are sorted
$median = $n/2;
foreach (@hkeys) {
   $median -= $hist{$_};
   if ($median <= 0) {
      $median = $_;
      last ;
   }
}

# and summarize
print "\n";
printf "items:    %i\n", $n;
printf "mean:     %g\n", $avg;
printf "mode:     %s\n", $mode;
printf "median:   %g\n", $median;
printf "stddev:   %g\n", $sdev;
