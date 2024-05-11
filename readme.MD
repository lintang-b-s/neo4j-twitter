

#### struktur node graph
```
Person { id: string}
Post { id: string}

```

#### struktur relasi/edge graph
```
1. relasi/edge follow
 (Person)-[r:FOLLOWS {since: follow.since}]->(otherUser: Person)

2. relasi/edge like
 (p:Person)-[r:LIKES {at: like.at}]->(likedPost:Post)
3. relasi/edge post 
    (p:Person)-[r:POSTED {at: post.at}]->(userPost:Post)
4. relasi/edge repost
    (p:Person)-[r:REPOSTS {at: repost.at}]->(repostedPost:Post)
```

### Quick Start
```
1. docker compose up -d 
2. docker logs neo4j -f (tunggu sampai neo4j started)
3. go mod tidy
4. go run main.go (tunggu sekitar 8-12 menit buat insert data ke neo4j)
```

#### Query ke neo4j
1. buka http://localhost:7474/browser/, login pake username=neo4j, password=lintang-neo4j
2. copy user_id buat query

##### Query siapa saja follower user
```

MATCH(followers: Person)-[:FOLLOWS]->(p: Person{id: "<user_id>"}) RETURN followers
```

#### Rekomendasi fanout tweets baru ke user mana aja
```
1. follower user
MATCH(followers: Person)-[:FOLLOWS]->(p: Person{id: "<user_id>"})
    RETURN followers

2.yang ngelike postingan user: 
WITH datetime()-duration({days: 210}) AS threshold
MATCH(orangYangLikePostsUser: Person)-[l:LIKES ]->(post)<-[:POSTED]-(p: Person{ id: "<user_id>"}) 
    WHERE l.at > threshold
    RETURN orangYangLikePostsUser

3. yang pernah repost postingan user:

WITH datetime()-duration({days: 210}) AS threshold
MATCH(orangYangRepostPostsUser: Person)-[l:REPOSTS ]->(post)<-[:POSTED]-(p: Person{ id: "<user_id>"}) 
    WHERE l.at > threshold
    RETURN orangYangRepostPostsUser


```



##### Query user follow siapa saja
```
MATCH(p: Person{id: "<user_id>"})-[:FOLLOWS]->(otherPerson) RETURN otherPerson
```

