import{S as Ve,i as Le,s as Ee,M as je,e as s,w as k,b as h,c as z,f as p,g as r,h as a,m as I,x as he,N as xe,P as Je,k as Ne,Q as Qe,n as ze,t as V,a as L,o as c,d as K,T as Ie,C as We,p as Ke,r as G,u as Ge}from"./index-a65ca895.js";import{S as Xe}from"./SdkTabs-ad912c8f.js";import{F as Ye}from"./FieldsQueryParam-ba250473.js";function Ue(i,l,o){const n=i.slice();return n[5]=l[o],n}function Be(i,l,o){const n=i.slice();return n[5]=l[o],n}function Fe(i,l){let o,n=l[5].code+"",m,g,u,b;function _(){return l[4](l[5])}return{key:i,first:null,c(){o=s("button"),m=k(n),g=h(),p(o,"class","tab-item"),G(o,"active",l[1]===l[5].code),this.first=o},m(v,A){r(v,o,A),a(o,m),a(o,g),u||(b=Ge(o,"click",_),u=!0)},p(v,A){l=v,A&4&&n!==(n=l[5].code+"")&&he(m,n),A&6&&G(o,"active",l[1]===l[5].code)},d(v){v&&c(o),u=!1,b()}}}function He(i,l){let o,n,m,g;return n=new je({props:{content:l[5].body}}),{key:i,first:null,c(){o=s("div"),z(n.$$.fragment),m=h(),p(o,"class","tab-item"),G(o,"active",l[1]===l[5].code),this.first=o},m(u,b){r(u,o,b),I(n,o,null),a(o,m),g=!0},p(u,b){l=u;const _={};b&4&&(_.content=l[5].body),n.$set(_),(!g||b&6)&&G(o,"active",l[1]===l[5].code)},i(u){g||(V(n.$$.fragment,u),g=!0)},o(u){L(n.$$.fragment,u),g=!1},d(u){u&&c(o),K(n)}}}function Ze(i){let l,o,n=i[0].name+"",m,g,u,b,_,v,A,P,X,S,E,pe,J,M,be,Y,N=i[0].name+"",Z,fe,ee,R,te,x,ae,W,le,y,oe,me,U,$,se,ge,ne,ke,f,_e,C,ve,we,Oe,ie,Ae,re,Se,ye,$e,ce,Te,Ce,q,ue,B,de,T,F,O=[],qe=new Map,De,H,w=[],Pe=new Map,D;v=new Xe({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${i[3]}');

        ...

        // OAuth2 authentication with a single realtime call.
        //
        // Make sure to register ${i[3]}/api/oauth2-redirect as redirect url.
        const authData = await pb.collection('users').authWithOAuth2({ provider: 'google' });

        // OR authenticate with manual OAuth2 code exchange
        // const authData = await pb.collection('users').authWithOAuth2Code(...);

        // after the above you can also access the auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.model.id);

        // "logout" the last authenticated model
        pb.authStore.clear();
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';
        import 'package:url_launcher/url_launcher.dart';

        final pb = PocketBase('${i[3]}');

        ...

        // OAuth2 authentication with a single realtime call.
        //
        // Make sure to register ${i[3]}/api/oauth2-redirect as redirect url.
        final authData = await pb.collection('users').authWithOAuth2('google', (url) async {
          await launchUrl(url);
        });

        // OR authenticate with manual OAuth2 code exchange
        // final authData = await pb.collection('users').authWithOAuth2Code(...);

        // after the above you can also access the auth data from the authStore
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.model.id);

        // "logout" the last authenticated model
        pb.authStore.clear();
    `}}),C=new je({props:{content:"?expand=relField1,relField2.subRelField"}}),q=new Ye({});let Q=i[2];const Me=e=>e[5].code;for(let e=0;e<Q.length;e+=1){let t=Be(i,Q,e),d=Me(t);qe.set(d,O[e]=Fe(d,t))}let j=i[2];const Re=e=>e[5].code;for(let e=0;e<j.length;e+=1){let t=Ue(i,j,e),d=Re(t);Pe.set(d,w[e]=He(d,t))}return{c(){l=s("h3"),o=k("Auth with OAuth2 ("),m=k(n),g=k(")"),u=h(),b=s("div"),b.innerHTML=`<p>Authenticate with an OAuth2 provider and returns a new auth token and record data.</p> 
    <p>For more details please check the
        <a href="https://pocketbase.io/docs/authentication/#oauth2-integration" target="_blank" rel="noopener noreferrer">OAuth2 integration documentation
        </a>.</p>`,_=h(),z(v.$$.fragment),A=h(),P=s("h6"),P.textContent="API details",X=h(),S=s("div"),E=s("strong"),E.textContent="POST",pe=h(),J=s("div"),M=s("p"),be=k("/api/collections/"),Y=s("strong"),Z=k(N),fe=k("/auth-with-oauth2"),ee=h(),R=s("div"),R.textContent="Body Parameters",te=h(),x=s("table"),x.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="50%">Description</th></tr></thead> 
    <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>provider</span></div></td> 
            <td><span class="label">String</span></td> 
            <td>The name of the OAuth2 client provider (eg. &quot;google&quot;).</td></tr> 
        <tr><td><div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>code</span></div></td> 
            <td><span class="label">String</span></td> 
            <td>The authorization code returned from the initial request.</td></tr> 
        <tr><td><div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>codeVerifier</span></div></td> 
            <td><span class="label">String</span></td> 
            <td>The code verifier sent with the initial request as part of the code_challenge.</td></tr> 
        <tr><td><div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>redirectUrl</span></div></td> 
            <td><span class="label">String</span></td> 
            <td>The redirect url sent with the initial request.</td></tr> 
        <tr><td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                    <span>createData</span></div></td> 
            <td><span class="label">Object</span></td> 
            <td><p>Optional data that will be used when creating the auth record on OAuth2 sign-up.</p> 
                <p>The created auth record must comply with the same requirements and validations in the
                    regular <strong>create</strong> action.
                    <br/> 
                    <em>The data can only be in <code>json</code>, aka. <code>multipart/form-data</code> and files
                        upload currently are not supported during OAuth2 sign-ups.</em></p></td></tr></tbody>`,ae=h(),W=s("div"),W.textContent="Query parameters",le=h(),y=s("table"),oe=s("thead"),oe.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr>`,me=h(),U=s("tbody"),$=s("tr"),se=s("td"),se.textContent="expand",ge=h(),ne=s("td"),ne.innerHTML='<span class="label">String</span>',ke=h(),f=s("td"),_e=k(`Auto expand record relations. Ex.:
                `),z(C.$$.fragment),ve=k(`
                Supports up to 6-levels depth nested relations expansion. `),we=s("br"),Oe=k(`
                The expanded relations will be appended to the record under the
                `),ie=s("code"),ie.textContent="expand",Ae=k(" property (eg. "),re=s("code"),re.textContent='"expand": {"relField1": {...}, ...}',Se=k(`).
                `),ye=s("br"),$e=k(`
                Only the relations to which the request user has permissions to `),ce=s("strong"),ce.textContent="view",Te=k(" will be expanded."),Ce=h(),z(q.$$.fragment),ue=h(),B=s("div"),B.textContent="Responses",de=h(),T=s("div"),F=s("div");for(let e=0;e<O.length;e+=1)O[e].c();De=h(),H=s("div");for(let e=0;e<w.length;e+=1)w[e].c();p(l,"class","m-b-sm"),p(b,"class","content txt-lg m-b-sm"),p(P,"class","m-b-xs"),p(E,"class","label label-primary"),p(J,"class","content"),p(S,"class","alert alert-success"),p(R,"class","section-title"),p(x,"class","table-compact table-border m-b-base"),p(W,"class","section-title"),p(y,"class","table-compact table-border m-b-base"),p(B,"class","section-title"),p(F,"class","tabs-header compact left"),p(H,"class","tabs-content"),p(T,"class","tabs")},m(e,t){r(e,l,t),a(l,o),a(l,m),a(l,g),r(e,u,t),r(e,b,t),r(e,_,t),I(v,e,t),r(e,A,t),r(e,P,t),r(e,X,t),r(e,S,t),a(S,E),a(S,pe),a(S,J),a(J,M),a(M,be),a(M,Y),a(Y,Z),a(M,fe),r(e,ee,t),r(e,R,t),r(e,te,t),r(e,x,t),r(e,ae,t),r(e,W,t),r(e,le,t),r(e,y,t),a(y,oe),a(y,me),a(y,U),a(U,$),a($,se),a($,ge),a($,ne),a($,ke),a($,f),a(f,_e),I(C,f,null),a(f,ve),a(f,we),a(f,Oe),a(f,ie),a(f,Ae),a(f,re),a(f,Se),a(f,ye),a(f,$e),a(f,ce),a(f,Te),a(U,Ce),I(q,U,null),r(e,ue,t),r(e,B,t),r(e,de,t),r(e,T,t),a(T,F);for(let d=0;d<O.length;d+=1)O[d]&&O[d].m(F,null);a(T,De),a(T,H);for(let d=0;d<w.length;d+=1)w[d]&&w[d].m(H,null);D=!0},p(e,[t]){(!D||t&1)&&n!==(n=e[0].name+"")&&he(m,n);const d={};t&8&&(d.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        // OAuth2 authentication with a single realtime call.
        //
        // Make sure to register ${e[3]}/api/oauth2-redirect as redirect url.
        const authData = await pb.collection('users').authWithOAuth2({ provider: 'google' });

        // OR authenticate with manual OAuth2 code exchange
        // const authData = await pb.collection('users').authWithOAuth2Code(...);

        // after the above you can also access the auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.model.id);

        // "logout" the last authenticated model
        pb.authStore.clear();
    `),t&8&&(d.dart=`
        import 'package:pocketbase/pocketbase.dart';
        import 'package:url_launcher/url_launcher.dart';

        final pb = PocketBase('${e[3]}');

        ...

        // OAuth2 authentication with a single realtime call.
        //
        // Make sure to register ${e[3]}/api/oauth2-redirect as redirect url.
        final authData = await pb.collection('users').authWithOAuth2('google', (url) async {
          await launchUrl(url);
        });

        // OR authenticate with manual OAuth2 code exchange
        // final authData = await pb.collection('users').authWithOAuth2Code(...);

        // after the above you can also access the auth data from the authStore
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.model.id);

        // "logout" the last authenticated model
        pb.authStore.clear();
    `),v.$set(d),(!D||t&1)&&N!==(N=e[0].name+"")&&he(Z,N),t&6&&(Q=e[2],O=xe(O,t,Me,1,e,Q,qe,F,Je,Fe,null,Be)),t&6&&(j=e[2],Ne(),w=xe(w,t,Re,1,e,j,Pe,H,Qe,He,null,Ue),ze())},i(e){if(!D){V(v.$$.fragment,e),V(C.$$.fragment,e),V(q.$$.fragment,e);for(let t=0;t<j.length;t+=1)V(w[t]);D=!0}},o(e){L(v.$$.fragment,e),L(C.$$.fragment,e),L(q.$$.fragment,e);for(let t=0;t<w.length;t+=1)L(w[t]);D=!1},d(e){e&&c(l),e&&c(u),e&&c(b),e&&c(_),K(v,e),e&&c(A),e&&c(P),e&&c(X),e&&c(S),e&&c(ee),e&&c(R),e&&c(te),e&&c(x),e&&c(ae),e&&c(W),e&&c(le),e&&c(y),K(C),K(q),e&&c(ue),e&&c(B),e&&c(de),e&&c(T);for(let t=0;t<O.length;t+=1)O[t].d();for(let t=0;t<w.length;t+=1)w[t].d()}}}function et(i,l,o){let n,{collection:m=new Ie}=l,g=200,u=[];const b=_=>o(1,g=_.code);return i.$$set=_=>{"collection"in _&&o(0,m=_.collection)},i.$$.update=()=>{i.$$.dirty&1&&o(2,u=[{code:200,body:JSON.stringify({token:"JWT_AUTH_TOKEN",record:We.dummyCollectionRecord(m),meta:{id:"abc123",name:"John Doe",username:"john.doe",email:"test@example.com",avatarUrl:"https://example.com/avatar.png",accessToken:"...",refreshToken:"...",rawUser:{}}},null,2)},{code:400,body:`
                {
                  "code": 400,
                  "message": "An error occurred while submitting the form.",
                  "data": {
                    "provider": {
                      "code": "validation_required",
                      "message": "Missing required value."
                    }
                  }
                }
            `}])},o(3,n=We.getApiExampleUrl(Ke.baseUrl)),[m,g,u,n,b]}class ot extends Ve{constructor(l){super(),Le(this,l,et,Ze,Ee,{collection:0})}}export{ot as default};
