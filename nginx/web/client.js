const {Subscription} = require('./src/pubnub_pb');
//const {PubSubClient} = require('./src/pubnub_grpc_web_pb');
const {DriverTrackingClient} = require('./src/pubnub.tracking_grpc_web_pb');
const {DriverTrackingMessage} = require('./src/pubnub.tracking_pb')

var echoService = new DriverTrackingClient('http://localhost:8080');
//var echoService = new DriverTrackingClient('pubnub-arke.prd-eks-bom-1.prd-eks.ps.pn:80')

//var message = new Message();
//message.setChannel('demo');
//message.setData({'test': true})
var subscription = new Subscription();
subscription.setChannel(process.env.CHANNEL)


//echoService.subscribe(subscription, {"publish-key": "pub-c-b9f19d65-c03f-41ad-a6b3-5ca8d2d25d75", "subscribe-key": "sub-c-deec24fc-903e-11e8-b601-f67fbeaec001"}, function(err, response) {
    //console.log(err);
    //console.log(response);
//});

var stream = echoService.subscribe(subscription, {"publish-key": process.env.PUBLISH_KEY, "subscribe-key": process.env.SUBSCRIBE_KEY});
stream.on('data', function(response) {
    //var msg = new LocationTrackingEnvelope(response);
  //console.log(msg.getData());
  var d = DriverTrackingMessage.toObject(true, response.getData());
  console.log("latitude=" + d.location.latitude);
  console.log("longitude=" + d.location.longitude);
  window.redraw(d.location.latitude, d.location.longitude);
  
});
stream.on('status', function(status) {
  console.log(status.code);
  console.log(status.details);
  console.log(status.metadata);
});
stream.on('end', function(end) {
  // stream end signal
});