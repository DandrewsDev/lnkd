import { writable } from "svelte/store";
import jwt_decode from "jwt-decode";


let storedJwt = writable(
    ""
);

if(typeof window !== "undefined") {
    storedJwt = writable(
        window.localStorage.getItem("userJwt") || ""
    );
    storedJwt.subscribe((val) => {
        window.localStorage.setItem("userJwt", val);
        if (val) {
            let decoded = jwt_decode(val);
            window.localStorage.setItem("username", decoded.username);
            window.localStorage.setItem("roles", decoded.roles);
        }
    });
}

export const userJwt = storedJwt