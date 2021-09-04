<template>
  <section>

    <h2>Rules</h2>

    <p>
      There are 4 levels of tree representing a tree's lifecycle:
    </p>
    <dl>
      <dt>Seed</dt>
      <dd>Seeds cannot cast seeds and do not gather Sun.</dd>

      <dt>Sprout</dt>
      <dd>Sprouts can cast seeds 1 space away and gather 1 Sun at the beginning of the Day.</dd>

      <dt>Sapling</dt>
      <dd>Saplings can cast seeds 2 spaces away and gather 2 Sun at the beginning of the Day.</dd>

      <dt>Tree</dt>
      <dd>Trees can cast seeds 3 spaces away and gather 3 Sun at the beginning of the Day. They also can be sold for
        points.
      </dd>
    </dl>

    <h2>Turns</h2>

    <p>
      Each player simultaneously chooses an action:
    </p>
    <ul>
      <li>
        Cast a seed from a tree
      </li>
      <li>
        Grow a tree
      </li>
      <li>
        Sell a tree
      </li>
      <li>
        End turn for the Day
      </li>
    </ul>

    <p>
      When a tree is used to perform an action (seeding, growing) it is dormant until the next day.
    </p>

    <p>
      Once both players End, the Sun direction moves and trees gather Energy.
    </p>

    <p>
      If one player Ends, the other may continue taking actions until they End too.
    </p>

    <p>After 26 days the game ends.</p>

    <h2>Costs</h2>

    <p>Actions cost ☀️ Energy.</p>

    <p>The base cost to cast a seed is 0.</p>

    <p>The base cost to grow a Seed is 1, a Sprout 3, and a Sapling 7.</p>

    <p>In addition to the base costs above, for each tree already at that size the cost increases by 1.</p>

    <p>To sell a tree costs 4.</p>

    <h2>Scoring</h2>

    <p>
      There is a global pool of nutrients which starts at 20. Each time a Tree is sold, the player that sells
      gains the points and the pool is reduced by 1.
    </p>

    <p>
      Bonus points are also gained if the Tree that is sold is on Rich (2) or Semi-Rich (1) soil.
    </p>

    <p>
      At the end of the game, players score 1 additional point for every 3 Energy.
    </p>

    <p>After 26 days the game ends, scores are tallied, and the highest score wins.</p>
  </section>
</template>
<style>
section {
  margin: 1rem;
}
</style>
