
    //Core API file that we will create 
var express = require('express');
var bodyParser = require('body-parser');
var app = express();
app.set("view engine", "pug");



var fs = require("fs");
app.use(bodyParser.json());

// Setting for Hyperledger Fabric

const { FileSystemWallet, Gateway } = require('fabric-network');
const path = require('path');
const ccpPath = path.resolve(__dirname, '..', '..', 'first-network', 'connection-org1.json');
app.set("views", path.join(__dirname, "views"));
console.log("Started API Server");


const renderListAsset = require("./views/viewsRenderingFunctions/listAsset");
const renderTrackAsset = require("./views/viewsRenderingFunctions/trackAsset");
var addToChain = require("./src/addAsset");



app.get('/', async function (req, res) {
        renderListAsset(req,res);
});


app.get('/trackAsset/:assetid', async function (req, res) {
       renderTrackAsset(req,res);
   
});



app.post('/addAsset', async function (req, res) {
       addToChain(req,res);
});


app.listen(8080);
