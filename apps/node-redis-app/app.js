const express = require("express");
const redis = require("redis");
const app = express();
const client = redis.createClient({ host: "redis-node", port: 6379 });

app.get("/health", (req, res) => {
  res.send("OK");
});

app.get("/", (req, res) => {
  client.get("visits", (err, visits) => {
    res.send("Number of visits is " + visits);
    client.set("visits", parseInt(visits) + 1);
  });
});

client.set("visits", 0);
app.listen(8080, () => {
  console.log("Listening on port 8080");
});
