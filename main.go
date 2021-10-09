package main

import (
	"context"
	"encoding/json"
	"fmt"
	// "log"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)
	
var client *mongo.Client

// db.users.createIndex( { "Name": 1 }, { unique: true } )

// indexName, err := coll.Indexes().CreateOne(
//     context.Background(),
//     mongo.IndexModel{
//         Keys:    bson.D{{Key: "Name", Value: 1}},
//         Options: options.Index().SetUnique(true),
//     },
// )

type Users struct {  //Users data structure
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
	Password  string             `json:"password,omitempty" bson:"password,omitempty"` //check type
}



type Posts struct {  //Posts data structure
	ID           primitive.ObjectID  `json:"_id,omitempty" bson:"_id,omitempty"`
	Caption      string              `json:"caption,omitempty" bson:"caption,omitempty"`
	Image_URL    string              `json:"image_url,omitempty" bson:"image_url,omitempty"`  //check type
	Timestamp    primitive.Timestamp `json:"timestamp,omitempty" bson:"timestamp,omitempty"`  //check type
}



func main() {
	fmt.Println("Backend API of Instagram...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("MONGO_URI")
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()
	router.HandleFunc("/users", CreateUsersEndpoint).Methods("POST")   
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")  //this needs to be changed a little
	router.HandleFunc("/users/{id}", GetUsersEndpoint).Methods("GET")
	router.HandleFunc("/posts", CreatePostsEndpoint).Methods("POST") 
	router.HandleFunc("/posts/{id}", GetPostsEndpoint).Methods("GET")
	http.ListenAndServe(":12345", router)
}

func CreateUsersEndpoint(response http.ResponseWriter, request *http.Request) {    // create an user- POST request
	response.Header().Set("content-type", "application/json")
	var person Users
	_ = json.NewDecoder(request.Body).Decode(&person)
	collection := client.Database("instagramPrototype").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, person)
	json.NewEncoder(response).Encode(result)
}


func GetUsersEndpoint(response http.ResponseWriter, request *http.Request) {     // get a user using id- GET request
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var person Users
	collection := client.Database("instagramPrototype").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Users{ID: id}).Decode(&person)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(person)
}



func GetUsersPostsEndpoint(response http.ResponseWriter, request *http.Request) {     // get a user using id- GET request
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var person Users
	var posts Posts
	collection := client.Database("instagramPrototype").Collection("users")
	collection2 := client.Database("instagramPrototype").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Users{ID: id}==Posts{ID: id}).Decode(&person)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(person)
}



func CreatePostsEndpoint(response http.ResponseWriter, request *http.Request) {    // create a post- POST request
	response.Header().Set("content-type", "application/json")
	var post Posts
	_ = json.NewDecoder(request.Body).Decode(&post)
	collection := client.Database("instagramPrototype").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, post)
	json.NewEncoder(response).Encode(result)
}

func GetPostsEndpoint(response http.ResponseWriter, request *http.Request) {     // get a post using id- GET request
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var post Posts
	collection := client.Database("instagramPrototype").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Posts{ID: id}).Decode(&post)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(post)
}


//this needs to be changed a little
func GetPeopleEndpoint(response http.ResponseWriter, request *http.Request) {  //GET all users, cursor to map through all the users
	response.Header().Set("content-type", "application/json")
	var people []Users
	collection := client.Database("instagramPrototype").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person Users
		cursor.Decode(&person)
		people = append(people, person)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(people)
}

//----
// func GetPostsOfUserEndpoint(response http.ResponseWriter, request *http.Request){
// 	collection := db.Collection("posts")
// 		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

// 		var posts []Posts

// 		cursor, _ := collection.Find(ctx, bson.M{})
// 		defer cursor.Close(ctx)

// 		for cursor.Next(ctx) {
// 			var post Posts
// 			cursor.Decode(&post)
// 			posts = append(posts, post)
// 		}

// 		return c.JSON(posts)
	
// }










// app.Get("/api/products/frontend", func(c *fiber.Ctx) error {
// 		collection := db.Collection("products")
// 		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

// 		var products []Product

// 		cursor, _ := collection.Find(ctx, bson.M{})
// 		defer cursor.Close(ctx)

// 		for cursor.Next(ctx) {
// 			var product Product
// 			cursor.Decode(&product)
// 			products = append(products, product)
// 		}

// 		return c.JSON(products)
// 	})




//----
// func GetPostsOfUserBackendEndpoint(response http.ResponseWriter, request *http.Request){
// collection := db.Collection("products")
// 		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

// 		var posts []Posts

// 		filter := bson.M{}
// 		findOptions := options.Find()

// 		if s := c.Query("s"); s != "" {
// 			filter = bson.M{
// 				"$or": []bson.M{
// 					{
// 						"title": bson.M{
// 							"$regex": primitive.Regex{
// 								Pattern: s,
// 								Options: "i",
// 							},
// 						},
// 					},
// 					{
// 						"description": bson.M{
// 							"$regex": primitive.Regex{
// 								Pattern: s,
// 								Options: "i",
// 							},
// 						},
// 					},
// 				},
// 			}

// }





// db.getCollection('users').aggregate([
//     {$match : {admin : 1}},
//     {$lookup: {from: "posts",localField: "_id",foreignField: "owner_id",as: "posts"}},
//     {$project : {
//             posts : { $filter : {input : "$posts"  , as : "post", cond : { $eq : ['$$post.via' , 'facebook'] } } },
//             admin : 1

//         }}

// ])



// db.getCollection('users').aggregate([
//     {$match : {admin : 1}},
//     {$lookup: {from: "posts",localField: "_id",foreignField: "owner_id",as: "posts"}},
//     {$unwind : "$posts"},
//     {$match : {"posts.via":"facebook"}},
//     { $group : {
//             _id : "$_id",
//             posts : {$push : "$posts"}
//     }}
// ])


// var myCursor = db.users.find({admin:1});
// var user_id = myCursor.hasNext() ? myCursor.next() : null;
// db.posts.find({owner_id : user_id._id});



// filterCursor, err := episodesCollection.Find(ctx, bson.M{"duration": 25})
// if err != nil {
//     log.Fatal(err)
// }
// var episodesFiltered []bson.M
// if err = filterCursor.All(ctx, &episodesFiltered); err != nil {
//     log.Fatal(err)
// }
// fmt.Println(episodesFiltered)











var mapUsers, mapComments, reduce;
db.users_comments.remove();

// setup sample data - wouldn't actually use this in production
db.users.remove();
db.posts.remove();
db.users.save({ID: ID[0]._id, Name:"Ross",email:"S@gmail.com",password:"password1"});
db.users.save({ID: ID[1]._id, Name:"Chandler" ,email:"Stm@gmail.com",password:"password2"});
db.users.save({ID: ID[2]._id, Name:"Monica",email:"mello@gmail.com",password:"password3"});

var users = db.users.find();
db.posts.save({ID: ID[0]._id, caption: "caption one", Image_URL="https://www.lovesove.com/wp-content/uploads/2018/07/har-bat-ke-jabab-me-muskura-hi-achchha-he-Attitude-Lines-in-hindi-Attitude-status-Attitude-quotes-lovesove.jpg" ,created: new ISODate()});
db.posts.save({ID: ID[1]._id, caption: "caption two", Image_URL="https://www.cleverfiles.com/howto/wp-content/uploads/2018/03/minion.jpg",created: new ISODate()});
db.posts.save({ID: ID[2]._id, caption: "caption three", Image_URL="https://cdn.pixabay.com/photo/2021/08/25/20/42/field-6574455__480.jpg",created: new ISODate()});
// end sample data setup

mapUsers = function() {
    var values = {
        Name: this.Name,
        email: this.email,
        password: this.password
    };
    emit(this._id, values);
};
mapComments = function() {
    var values = {
        commentId: this._id,
        comment: this.comment,
        created: this.created
    };
    emit(this.userId, values);
};
reduce = function(k, values) {
    var result = {}, commentFields = {
        "commentId": '', 
        "comment": '',
        "created": ''
    };
    values.forEach(function(value) {
        var field;
        if ("comment" in value) {
            if (!("comments" in result)) {
                result.comments = [];
            }
            result.comments.push(value);
        } else if ("comments" in value) {
            if (!("comments" in result)) {
                result.comments = [];
            }
            result.comments.push.apply(result.comments, value.comments);
        }
        for (field in value) {
            if (value.hasOwnProperty(field) && !(field in commentFields)) {
                result[field] = value[field];
            }
        }
    });
    return result;
};
db.users.mapReduce(mapUsers, reduce, {"out": {"reduce": "users_comments"}});
db.comments.mapReduce(mapComments, reduce, {"out": {"reduce": "users_comments"}});
db.users_comments.find().pretty(); // see the resulting collection