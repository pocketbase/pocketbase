import{S as ee,i as te,s as le,T as U,e as p,b as S,w as V,f as g,g as b,h,x as z,U as G,V as ge,k as Z,W as ke,n as x,t as L,a as N,o as _,C as oe,r as Y,u as ae,R as we,c as K,m as Q,d as X,Q as $e,X as se,p as Te,Y as ne}from"./index-ktUKDHmL.js";function ie(s,t,e){const a=s.slice();return a[4]=t[e],a}function ce(s,t,e){const a=s.slice();return a[4]=t[e],a}function re(s,t){let e,a=t[4].code+"",d,c,r,n;function u(){return t[3](t[4])}return{key:s,first:null,c(){e=p("button"),d=V(a),c=S(),g(e,"class","tab-item"),Y(e,"active",t[1]===t[4].code),this.first=e},m(m,q){b(m,e,q),h(e,d),h(e,c),r||(n=ae(e,"click",u),r=!0)},p(m,q){t=m,q&4&&a!==(a=t[4].code+"")&&z(d,a),q&6&&Y(e,"active",t[1]===t[4].code)},d(m){m&&_(e),r=!1,n()}}}function de(s,t){let e,a,d,c;return a=new we({props:{content:t[4].body}}),{key:s,first:null,c(){e=p("div"),K(a.$$.fragment),d=S(),g(e,"class","tab-item"),Y(e,"active",t[1]===t[4].code),this.first=e},m(r,n){b(r,e,n),Q(a,e,null),h(e,d),c=!0},p(r,n){t=r;const u={};n&4&&(u.content=t[4].body),a.$set(u),(!c||n&6)&&Y(e,"active",t[1]===t[4].code)},i(r){c||(L(a.$$.fragment,r),c=!0)},o(r){N(a.$$.fragment,r),c=!1},d(r){r&&_(e),X(a)}}}function Pe(s){let t,e,a,d,c,r,n,u=s[0].name+"",m,q,M,C,B,A,H,R,W,y,P,w=[],$=new Map,E,D,k=[],J=new Map,I,i=U(s[2]);const v=l=>l[4].code;for(let l=0;l<i.length;l+=1){let o=ce(s,i,l),f=v(o);$.set(f,w[l]=re(f,o))}let O=U(s[2]);const j=l=>l[4].code;for(let l=0;l<O.length;l+=1){let o=ie(s,O,l),f=j(o);J.set(f,k[l]=de(f,o))}return{c(){t=p("div"),e=p("strong"),e.textContent="POST",a=S(),d=p("div"),c=p("p"),r=V("/api/collections/"),n=p("strong"),m=V(u),q=V("/auth-with-otp"),M=S(),C=p("div"),C.textContent="Body Parameters",B=S(),A=p("table"),A.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr></thead> <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>otpId</span></div></td> <td><span class="label">String</span></td> <td>The id of the OTP request.</td></tr> <tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>password</span></div></td> <td><span class="label">String</span></td> <td>The one-time password.</td></tr></tbody>',H=S(),R=p("div"),R.textContent="Responses",W=S(),y=p("div"),P=p("div");for(let l=0;l<w.length;l+=1)w[l].c();E=S(),D=p("div");for(let l=0;l<k.length;l+=1)k[l].c();g(e,"class","label label-primary"),g(d,"class","content"),g(t,"class","alert alert-success"),g(C,"class","section-title"),g(A,"class","table-compact table-border m-b-base"),g(R,"class","section-title"),g(P,"class","tabs-header compact combined left"),g(D,"class","tabs-content"),g(y,"class","tabs")},m(l,o){b(l,t,o),h(t,e),h(t,a),h(t,d),h(d,c),h(c,r),h(c,n),h(n,m),h(c,q),b(l,M,o),b(l,C,o),b(l,B,o),b(l,A,o),b(l,H,o),b(l,R,o),b(l,W,o),b(l,y,o),h(y,P);for(let f=0;f<w.length;f+=1)w[f]&&w[f].m(P,null);h(y,E),h(y,D);for(let f=0;f<k.length;f+=1)k[f]&&k[f].m(D,null);I=!0},p(l,[o]){(!I||o&1)&&u!==(u=l[0].name+"")&&z(m,u),o&6&&(i=U(l[2]),w=G(w,o,v,1,l,i,$,P,ge,re,null,ce)),o&6&&(O=U(l[2]),Z(),k=G(k,o,j,1,l,O,J,D,ke,de,null,ie),x())},i(l){if(!I){for(let o=0;o<O.length;o+=1)L(k[o]);I=!0}},o(l){for(let o=0;o<k.length;o+=1)N(k[o]);I=!1},d(l){l&&(_(t),_(M),_(C),_(B),_(A),_(H),_(R),_(W),_(y));for(let o=0;o<w.length;o+=1)w[o].d();for(let o=0;o<k.length;o+=1)k[o].d()}}}function Oe(s,t,e){let{collection:a}=t,d=200,c=[];const r=n=>e(1,d=n.code);return s.$$set=n=>{"collection"in n&&e(0,a=n.collection)},s.$$.update=()=>{s.$$.dirty&1&&e(2,c=[{code:200,body:JSON.stringify({token:"JWT_TOKEN",record:oe.dummyCollectionRecord(a)},null,2)},{code:400,body:`
                {
                  "code": 400,
                  "message": "Failed to authenticate.",
                  "data": {
                    "otpId": {
                      "code": "validation_required",
                      "message": "Missing required value."
                    }
                  }
                }
            `}])},[a,d,c,r]}class Se extends ee{constructor(t){super(),te(this,t,Oe,Pe,le,{collection:0})}}function ue(s,t,e){const a=s.slice();return a[4]=t[e],a}function he(s,t,e){const a=s.slice();return a[4]=t[e],a}function pe(s,t){let e,a=t[4].code+"",d,c,r,n;function u(){return t[3](t[4])}return{key:s,first:null,c(){e=p("button"),d=V(a),c=S(),g(e,"class","tab-item"),Y(e,"active",t[1]===t[4].code),this.first=e},m(m,q){b(m,e,q),h(e,d),h(e,c),r||(n=ae(e,"click",u),r=!0)},p(m,q){t=m,q&4&&a!==(a=t[4].code+"")&&z(d,a),q&6&&Y(e,"active",t[1]===t[4].code)},d(m){m&&_(e),r=!1,n()}}}function fe(s,t){let e,a,d,c;return a=new we({props:{content:t[4].body}}),{key:s,first:null,c(){e=p("div"),K(a.$$.fragment),d=S(),g(e,"class","tab-item"),Y(e,"active",t[1]===t[4].code),this.first=e},m(r,n){b(r,e,n),Q(a,e,null),h(e,d),c=!0},p(r,n){t=r;const u={};n&4&&(u.content=t[4].body),a.$set(u),(!c||n&6)&&Y(e,"active",t[1]===t[4].code)},i(r){c||(L(a.$$.fragment,r),c=!0)},o(r){N(a.$$.fragment,r),c=!1},d(r){r&&_(e),X(a)}}}function ye(s){let t,e,a,d,c,r,n,u=s[0].name+"",m,q,M,C,B,A,H,R,W,y,P,w=[],$=new Map,E,D,k=[],J=new Map,I,i=U(s[2]);const v=l=>l[4].code;for(let l=0;l<i.length;l+=1){let o=he(s,i,l),f=v(o);$.set(f,w[l]=pe(f,o))}let O=U(s[2]);const j=l=>l[4].code;for(let l=0;l<O.length;l+=1){let o=ue(s,O,l),f=j(o);J.set(f,k[l]=fe(f,o))}return{c(){t=p("div"),e=p("strong"),e.textContent="POST",a=S(),d=p("div"),c=p("p"),r=V("/api/collections/"),n=p("strong"),m=V(u),q=V("/request-otp"),M=S(),C=p("div"),C.textContent="Body Parameters",B=S(),A=p("table"),A.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr></thead> <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>email</span></div></td> <td><span class="label">String</span></td> <td>The auth record email address to send the OTP request (if exists).</td></tr></tbody>',H=S(),R=p("div"),R.textContent="Responses",W=S(),y=p("div"),P=p("div");for(let l=0;l<w.length;l+=1)w[l].c();E=S(),D=p("div");for(let l=0;l<k.length;l+=1)k[l].c();g(e,"class","label label-primary"),g(d,"class","content"),g(t,"class","alert alert-success"),g(C,"class","section-title"),g(A,"class","table-compact table-border m-b-base"),g(R,"class","section-title"),g(P,"class","tabs-header compact combined left"),g(D,"class","tabs-content"),g(y,"class","tabs")},m(l,o){b(l,t,o),h(t,e),h(t,a),h(t,d),h(d,c),h(c,r),h(c,n),h(n,m),h(c,q),b(l,M,o),b(l,C,o),b(l,B,o),b(l,A,o),b(l,H,o),b(l,R,o),b(l,W,o),b(l,y,o),h(y,P);for(let f=0;f<w.length;f+=1)w[f]&&w[f].m(P,null);h(y,E),h(y,D);for(let f=0;f<k.length;f+=1)k[f]&&k[f].m(D,null);I=!0},p(l,[o]){(!I||o&1)&&u!==(u=l[0].name+"")&&z(m,u),o&6&&(i=U(l[2]),w=G(w,o,v,1,l,i,$,P,ge,pe,null,he)),o&6&&(O=U(l[2]),Z(),k=G(k,o,j,1,l,O,J,D,ke,fe,null,ue),x())},i(l){if(!I){for(let o=0;o<O.length;o+=1)L(k[o]);I=!0}},o(l){for(let o=0;o<k.length;o+=1)N(k[o]);I=!1},d(l){l&&(_(t),_(M),_(C),_(B),_(A),_(H),_(R),_(W),_(y));for(let o=0;o<w.length;o+=1)w[o].d();for(let o=0;o<k.length;o+=1)k[o].d()}}}function qe(s,t,e){let{collection:a}=t,d=200,c=[];const r=n=>e(1,d=n.code);return s.$$set=n=>{"collection"in n&&e(0,a=n.collection)},e(2,c=[{code:200,body:JSON.stringify({otpId:oe.randomString(15)},null,2)},{code:400,body:`
                {
                  "code": 400,
                  "message": "An error occurred while validating the submitted data.",
                  "data": {
                    "email": {
                      "code": "validation_is_email",
                      "message": "Must be a valid email address."
                    }
                  }
                }
            `},{code:429,body:`
                {
                  "code": 429,
                  "message": "You've send too many OTP requests, please try again later.",
                  "data": {}
                }
            `}]),[a,d,c,r]}class Ae extends ee{constructor(t){super(),te(this,t,qe,ye,le,{collection:0})}}function me(s,t,e){const a=s.slice();return a[5]=t[e],a[7]=e,a}function be(s,t,e){const a=s.slice();return a[5]=t[e],a[7]=e,a}function _e(s){let t,e,a,d,c;function r(){return s[4](s[7])}return{c(){t=p("button"),e=p("div"),e.textContent=`${s[5].title}`,a=S(),g(e,"class","txt"),g(t,"class","tab-item"),Y(t,"active",s[1]==s[7])},m(n,u){b(n,t,u),h(t,e),h(t,a),d||(c=ae(t,"click",r),d=!0)},p(n,u){s=n,u&2&&Y(t,"active",s[1]==s[7])},d(n){n&&_(t),d=!1,c()}}}function ve(s){let t,e,a,d;var c=s[5].component;function r(n,u){return{props:{collection:n[0]}}}return c&&(e=ne(c,r(s))),{c(){t=p("div"),e&&K(e.$$.fragment),a=S(),g(t,"class","tab-item"),Y(t,"active",s[1]==s[7])},m(n,u){b(n,t,u),e&&Q(e,t,null),h(t,a),d=!0},p(n,u){if(c!==(c=n[5].component)){if(e){Z();const m=e;N(m.$$.fragment,1,0,()=>{X(m,1)}),x()}c?(e=ne(c,r(n)),K(e.$$.fragment),L(e.$$.fragment,1),Q(e,t,a)):e=null}else if(c){const m={};u&1&&(m.collection=n[0]),e.$set(m)}(!d||u&2)&&Y(t,"active",n[1]==n[7])},i(n){d||(e&&L(e.$$.fragment,n),d=!0)},o(n){e&&N(e.$$.fragment,n),d=!1},d(n){n&&_(t),e&&X(e)}}}function Re(s){var D,k,J,I;let t,e,a=s[0].name+"",d,c,r,n,u,m,q,M,C,B,A,H,R,W;m=new $e({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${s[2]}');

        ...

        // send OTP email to the provided auth record
        const req = await pb.collection('${(D=s[0])==null?void 0:D.name}').requestOTP('test@example.com');

        // ... show a screen/popup to enter the password from the email ...

        // authenticate with the requested OTP id and the email password
        const authData = await pb.collection('${(k=s[0])==null?void 0:k.name}').authWithOTP(
            req.otpId,
            "YOUR_OTP",
        );

        // after the above you can also access the auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.record.id);

        // "logout"
        pb.authStore.clear();
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${s[2]}');

        ...

        // send OTP email to the provided auth record
        final req = await pb.collection('${(J=s[0])==null?void 0:J.name}').requestOTP('test@example.com');

        // ... show a screen/popup to enter the password from the email ...

        // authenticate with the requested OTP id and the email password
        final authData = await pb.collection('${(I=s[0])==null?void 0:I.name}').authWithOTP(
            req.otpId,
            "YOUR_OTP",
        );

        // after the above you can also access the auth data from the authStore
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.record.id);

        // "logout"
        pb.authStore.clear();
    `}});let y=U(s[3]),P=[];for(let i=0;i<y.length;i+=1)P[i]=_e(be(s,y,i));let w=U(s[3]),$=[];for(let i=0;i<w.length;i+=1)$[i]=ve(me(s,w,i));const E=i=>N($[i],1,1,()=>{$[i]=null});return{c(){t=p("h3"),e=V("Auth with OTP ("),d=V(a),c=V(")"),r=S(),n=p("div"),n.innerHTML=`<p>Authenticate with an one-time password (OTP).</p> <p>Note that when requesting an OTP we return an <code>otpId</code> even if a user with the provided email
        doesn&#39;t exist as a very basic enumeration protection.</p>`,u=S(),K(m.$$.fragment),q=S(),M=p("h6"),M.textContent="API details",C=S(),B=p("div"),A=p("div");for(let i=0;i<P.length;i+=1)P[i].c();H=S(),R=p("div");for(let i=0;i<$.length;i+=1)$[i].c();g(t,"class","m-b-sm"),g(n,"class","content txt-lg m-b-sm"),g(M,"class","m-b-xs"),g(A,"class","tabs-header compact"),g(R,"class","tabs-content"),g(B,"class","tabs")},m(i,v){b(i,t,v),h(t,e),h(t,d),h(t,c),b(i,r,v),b(i,n,v),b(i,u,v),Q(m,i,v),b(i,q,v),b(i,M,v),b(i,C,v),b(i,B,v),h(B,A);for(let O=0;O<P.length;O+=1)P[O]&&P[O].m(A,null);h(B,H),h(B,R);for(let O=0;O<$.length;O+=1)$[O]&&$[O].m(R,null);W=!0},p(i,[v]){var j,l,o,f;(!W||v&1)&&a!==(a=i[0].name+"")&&z(d,a);const O={};if(v&5&&(O.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${i[2]}');

        ...

        // send OTP email to the provided auth record
        const req = await pb.collection('${(j=i[0])==null?void 0:j.name}').requestOTP('test@example.com');

        // ... show a screen/popup to enter the password from the email ...

        // authenticate with the requested OTP id and the email password
        const authData = await pb.collection('${(l=i[0])==null?void 0:l.name}').authWithOTP(
            req.otpId,
            "YOUR_OTP",
        );

        // after the above you can also access the auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.record.id);

        // "logout"
        pb.authStore.clear();
    `),v&5&&(O.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${i[2]}');

        ...

        // send OTP email to the provided auth record
        final req = await pb.collection('${(o=i[0])==null?void 0:o.name}').requestOTP('test@example.com');

        // ... show a screen/popup to enter the password from the email ...

        // authenticate with the requested OTP id and the email password
        final authData = await pb.collection('${(f=i[0])==null?void 0:f.name}').authWithOTP(
            req.otpId,
            "YOUR_OTP",
        );

        // after the above you can also access the auth data from the authStore
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.record.id);

        // "logout"
        pb.authStore.clear();
    `),m.$set(O),v&10){y=U(i[3]);let T;for(T=0;T<y.length;T+=1){const F=be(i,y,T);P[T]?P[T].p(F,v):(P[T]=_e(F),P[T].c(),P[T].m(A,null))}for(;T<P.length;T+=1)P[T].d(1);P.length=y.length}if(v&11){w=U(i[3]);let T;for(T=0;T<w.length;T+=1){const F=me(i,w,T);$[T]?($[T].p(F,v),L($[T],1)):($[T]=ve(F),$[T].c(),L($[T],1),$[T].m(R,null))}for(Z(),T=w.length;T<$.length;T+=1)E(T);x()}},i(i){if(!W){L(m.$$.fragment,i);for(let v=0;v<w.length;v+=1)L($[v]);W=!0}},o(i){N(m.$$.fragment,i),$=$.filter(Boolean);for(let v=0;v<$.length;v+=1)N($[v]);W=!1},d(i){i&&(_(t),_(r),_(n),_(u),_(q),_(M),_(C),_(B)),X(m,i),se(P,i),se($,i)}}}function Ce(s,t,e){let a,{collection:d}=t;const c=[{title:"OTP Request",component:Ae},{title:"OTP Auth",component:Se}];let r=0;const n=u=>e(1,r=u);return s.$$set=u=>{"collection"in u&&e(0,d=u.collection)},e(2,a=oe.getApiExampleUrl(Te.baseURL)),[d,r,a,c,n]}class De extends ee{constructor(t){super(),te(this,t,Ce,Re,le,{collection:0})}}export{De as default};
