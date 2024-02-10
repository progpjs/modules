import "@progp/core"
import {HttpRequest, HttpServer} from "@progp/http"

let server = new HttpServer(8000);
let host = server.getHost("localhost");

let count = 0;

async function handler(req: HttpRequest): Promise<void> {
    //console.log("Request IP:" + req.requestIP());
    //console.log("Request path: ", req.requestPath());

    req.returnHtml(200, "call " + count++)

}

host.GET("/", handler);

server.start();