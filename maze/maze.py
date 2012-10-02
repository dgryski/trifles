#!/usr/bin/python

# WARNING: Some of this code is from 2004.  Let's just pretend all the ugly
# parts are the old stuff and all the cool/interesting/good pieces are from my
# 2012 rework :)

import random

class Cell:
   def __init__(self, xy):
      self.u = 1
      self.d = 1
      self.r = 1
      self.l = 1
      self.visited = 0
      self.xy = xy
      self.solnpath = 0

def remove_idx_move(idx, nodes):
   n = nodes.pop()
   if idx != len(nodes):
       nodes[idx] = n

def remove_idx(idx, nodes):
    nodes.pop(idx)

def most_recent_node(nodes):
    return len(nodes)-1

def random_node(nodes):
    return random.randrange(0, len(nodes))

def random_recent_node(how_recent):
    return lambda nodes:len(nodes) - 1 - random.randrange(0, min(how_recent, len(nodes)))

def random_old_node(how_old):
    return lambda nodes:random.randrange(0, min(how_old, len(nodes)))

class Maze:
   def __init__(self, height, width):
      self.walls = []
      for i in range(height):
         b = []
         for j in range(width): b.append(Cell((j,i)))
         self.walls.append(b)

      self.height = height
      self.width = width
      self.entrance = (0,0)
      self.exit = (height-1, width-1)
      self.soln = None

   # by providing choice_strategy and removal_strategy, we can construct different maze types
   def growing_tree(self, choice_strategy, removal_strategy):
      nodes = [self.walls[0][0]]

      while nodes:

         idx = choice_strategy(nodes)

         curr = nodes[idx]
         (x,y) = curr.xy

         # make a list of unvisited neighbours
         n = []
         if (x > 0): n.append(self.walls[y][x-1]);
         if (x < self.width-1): n.append(self.walls[y][x+1]);
         if (y > 0): n.append(self.walls[y-1][x]);
         if (y < self.height-1): n.append(self.walls[y+1][x]);

         n = filter(lambda x: x.visited == 0, n)

         if n:
            # pick random neighbour
            randnode = random.choice(n)
            (i,j) = randnode.xy

            randnode.visited = 1

            # knock down wall between (x,y) and (i,j)
            if   (i < x):
                self.walls[y][x].l = 0
                self.walls[j][i].r = 0
            elif (i > x):
                self.walls[y][x].r = 0
                self.walls[j][i].l = 0
            elif (j < y):
                self.walls[y][x].u = 0
                self.walls[j][i].d = 0
            elif (j > y):
                self.walls[y][x].d = 0
                self.walls[j][i].u = 0

            nodes.append(randnode)
         else:
            # we're done with this node -- remove from list
            removal_strategy(idx, nodes)

   def recursive_backtracker(self):
       self.growing_tree(most_recent_node, remove_idx)

   def nearly_prims(self):
       self.growing_tree(random_node, remove_idx)

   def low_river(self, how_recent):
       self.growing_tree(random_recent_node(how_recent), remove_idx)

   def solve(self):
      for i in range(self.height):
         for j in range(self.width): self.walls[i][j].visited = 0

      start = self.walls[self.entrance[0]][self.entrance[1]]
      exit = self.walls[self.exit[0]][self.exit[1]]

      start.visited = 1
      nodes = [start]

      while nodes:
         curr = nodes[-1]

         if curr.xy[0] == exit.xy[0] and curr.xy[1] == exit.xy[1]:
            break

         (x,y) = curr.xy

         # make a list of unvisited neighbours
         n = []
         if not curr.d: n.append(self.walls[y+1][x])
         if not curr.u: n.append(self.walls[y-1][x])
         if not curr.l: n.append(self.walls[y][x-1])
         if not curr.r: n.append(self.walls[y][x+1])
         n = filter(lambda x: x.visited == 0, n)

         if n:
            n[0].visited = 1
            nodes.append(n[0])
         else:
             nodes.pop()

      for n in nodes:
         n.solnpath = 1

if __name__ == '__main__':
   maze = Maze(20,20)

#   maze.recursive_backtracker()
#   maze.nearly_prims()
   maze.low_river(5)

   def draw_cell(x):
      if x.d: d = '_'
      else: d = ' '
      if x.r: r = '|'
      else: r = '_'
      return '%c%c' % (d,r)

   def draw_bottom_cell(x):
      if x.r: r = '|'
      else: r = '_'
      return '_%c' % r

   def draw_wall(x):
      p = ' '
      if x.solnpath:
          p = '\033[91m' + '*' + '\033[0m'
      s = ''
      if x.r:
          return p+'#'
      return p+p

   def draw_floor(x):
      p = ' '
      if x.solnpath:
          p = '\033[91m' + '*' + '\033[0m'
      if x.d: return '##'
      else: return p+'#'

   maze.solve()

   print '_' * (maze.width*2+1)
   for i in maze.walls[:-1]:
      print '|%s' % ''.join(map(draw_cell, i))
   print '|%s' % ''.join(map(draw_bottom_cell, maze.walls[-1]))

   print 

   print '#' * (maze.width*2+1)
   for i in maze.walls:
      print '#%s' % ''.join(map(draw_wall, i))
      print '#%s' % ''.join(map(draw_floor, i))
