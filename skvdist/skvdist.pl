#!/usr/bin/perl -w


use strict;

use ShardedKV::Continuum::Jump;

my $continuum = ShardedKV::Continuum::Jump->new(
    from => {
        ids => [map "shard$_" , 1..32],
    }
);

my $processed = 0;
while(<>) {
    chomp;
    my ($shard, $key) = split / /, $_, 2;
    my $dst = $continuum->choose($key);
    die "mismatch for $key: $shard != $dst" if $shard ne $dst;
    $processed++;
    if (($processed % (1024*1024)) == 0) {
        print scalar localtime, " ", $processed, "\n";
    }
}

print "processed $processed files\n";
