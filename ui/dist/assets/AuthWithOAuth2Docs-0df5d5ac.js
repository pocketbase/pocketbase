import{S as Le,i as Ee,s as Je,M as Ve,N as z,e as o,w as k,b as h,c as I,f as p,g as r,h as a,m as K,x as pe,P as Ue,Q as Ne,k as Qe,R as ze,n as Ie,t as V,a as L,o as c,d as G,U as Ke,C as We,p as Ge,r as X,u as Xe}from"./index-a084d9d7.js";import{S as Ye}from"./SdkTabs-ba0ec979.js";import{F as Ze}from"./FieldsQueryParam-71e01e64.js";function Be(i,l,s){const n=i.slice();return n[5]=l[s],n}function Fe(i,l,s){const n=i.slice();return n[5]=l[s],n}function He(i,l){let s,n=l[5].code+"",f,g,d,b;function _(){return l[4](l[5])}return{key:i,first:null,c(){s=o("button"),f=k(n),g=h(),p(s,"class","tab-item"),X(s,"active",l[1]===l[5].code),this.first=s},m(v,A){r(v,s,A),a(s,f),a(s,g),d||(b=Xe(s,"click",_),d=!0)},p(v,A){l=v,A&4&&n!==(n=l[5].code+"")&&pe(f,n),A&6&&X(s,"active",l[1]===l[5].code)},d(v){v&&c(s),d=!1,b()}}}function je(i,l){let s,n,f,g;return n=new Ve({props:{content:l[5].body}}),{key:i,first:null,c(){s=o("div"),I(n.$$.fragment),f=h(),p(s,"class","tab-item"),X(s,"active",l[1]===l[5].code),this.first=s},m(d,b){r(d,s,b),K(n,s,null),a(s,f),g=!0},p(d,b){l=d;const _={};b&4&&(_.content=l[5].body),n.$set(_),(!g||b&6)&&X(s,"active",l[1]===l[5].code)},i(d){g||(V(n.$$.fragment,d),g=!0)},o(d){L(n.$$.fragment,d),g=!1},d(d){d&&c(s),G(n)}}}function et(i){let l,s,n=i[0].name+"",f,g,d,b,_,v,A,P,Y,S,E,be,J,M,me,Z,N=i[0].name+"",ee,fe,te,R,ae,x,le,U,se,y,oe,ge,W,$,ne,ke,ie,_e,m,ve,C,we,Oe,Ae,re,Se,ce,ye,$e,Te,de,Ce,qe,q,ue,B,he,T,F,O=[],De=new Map,Pe,H,w=[],Me=new Map,D;v=new Ye({props:{js:`
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
    `}}),C=new Ve({props:{content:"?expand=relField1,relField2.subRelField"}}),q=new Ze({});let Q=z(i[2]);const Re=e=>e[5].code;for(let e=0;e<Q.length;e+=1){let t=Fe(i,Q,e),u=Re(t);De.set(u,O[e]=He(u,t))}let j=z(i[2]);const xe=e=>e[5].code;for(let e=0;e<j.length;e+=1){let t=Be(i,j,e),u=xe(t);Me.set(u,w[e]=je(u,t))}return{c(){l=o("h3"),s=k("Auth with OAuth2 ("),f=k(n),g=k(")"),d=h(),b=o("div"),b.innerHTML=`<p>Authenticate with an OAuth2 provider and returns a new auth token and record data.</p> <p>For more details please check the
        <a href="https://pocketbase.io/docs/authentication/#oauth2-integration" target="_blank" rel="noopener noreferrer">OAuth2 integration documentation
        </a>.</p>`,_=h(),I(v.$$.fragment),A=h(),P=o("h6"),P.textContent="API details",Y=h(),S=o("div"),E=o("strong"),E.textContent="POST",be=h(),J=o("div"),M=o("p"),me=k("/api/collections/"),Z=o("strong"),ee=k(N),fe=k("/auth-with-oauth2"),te=h(),R=o("div"),R.textContent="Body Parameters",ae=h(),x=o("table"),x.innerHTML=`<thead><tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr></thead> <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>provider</span></div></td> <td><span class="label">String</span></td> <td>The name of the OAuth2 client provider (eg. &quot;google&quot;).</td></tr> <tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>code</span></div></td> <td><span class="label">String</span></td> <td>The authorization code returned from the initial request.</td></tr> <tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>codeVerifier</span></div></td> <td><span class="label">String</span></td> <td>The code verifier sent with the initial request as part of the code_challenge.</td></tr> <tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>redirectUrl</span></div></td> <td><span class="label">String</span></td> <td>The redirect url sent with the initial request.</td></tr> <tr><td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>createData</span></div></td> <td><span class="label">Object</span></td> <td><p>Optional data that will be used when creating the auth record on OAuth2 sign-up.</p> <p>The created auth record must comply with the same requirements and validations in the
                    regular <strong>create</strong> action.
                    <br/> <em>The data can only be in <code>json</code>, aka. <code>multipart/form-data</code> and files
                        upload currently are not supported during OAuth2 sign-ups.</em></p></td></tr></tbody>`,le=h(),U=o("div"),U.textContent="Query parameters",se=h(),y=o("table"),oe=o("thead"),oe.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',ge=h(),W=o("tbody"),$=o("tr"),ne=o("td"),ne.textContent="expand",ke=h(),ie=o("td"),ie.innerHTML='<span class="label">String</span>',_e=h(),m=o("td"),ve=k(`Auto expand record relations. Ex.:
                `),I(C.$$.fragment),we=k(`
                Supports up to 6-levels depth nested relations expansion. `),Oe=o("br"),Ae=k(`
                The expanded relations will be appended to the record under the
                `),re=o("code"),re.textContent="expand",Se=k(" property (eg. "),ce=o("code"),ce.textContent='"expand": {"relField1": {...}, ...}',ye=k(`).
                `),$e=o("br"),Te=k(`
                Only the relations to which the request user has permissions to `),de=o("strong"),de.textContent="view",Ce=k(" will be expanded."),qe=h(),I(q.$$.fragment),ue=h(),B=o("div"),B.textContent="Responses",he=h(),T=o("div"),F=o("div");for(let e=0;e<O.length;e+=1)O[e].c();Pe=h(),H=o("div");for(let e=0;e<w.length;e+=1)w[e].c();p(l,"class","m-b-sm"),p(b,"class","content txt-lg m-b-sm"),p(P,"class","m-b-xs"),p(E,"class","label label-primary"),p(J,"class","content"),p(S,"class","alert alert-success"),p(R,"class","section-title"),p(x,"class","table-compact table-border m-b-base"),p(U,"class","section-title"),p(y,"class","table-compact table-border m-b-base"),p(B,"class","section-title"),p(F,"class","tabs-header compact left"),p(H,"class","tabs-content"),p(T,"class","tabs")},m(e,t){r(e,l,t),a(l,s),a(l,f),a(l,g),r(e,d,t),r(e,b,t),r(e,_,t),K(v,e,t),r(e,A,t),r(e,P,t),r(e,Y,t),r(e,S,t),a(S,E),a(S,be),a(S,J),a(J,M),a(M,me),a(M,Z),a(Z,ee),a(M,fe),r(e,te,t),r(e,R,t),r(e,ae,t),r(e,x,t),r(e,le,t),r(e,U,t),r(e,se,t),r(e,y,t),a(y,oe),a(y,ge),a(y,W),a(W,$),a($,ne),a($,ke),a($,ie),a($,_e),a($,m),a(m,ve),K(C,m,null),a(m,we),a(m,Oe),a(m,Ae),a(m,re),a(m,Se),a(m,ce),a(m,ye),a(m,$e),a(m,Te),a(m,de),a(m,Ce),a(W,qe),K(q,W,null),r(e,ue,t),r(e,B,t),r(e,he,t),r(e,T,t),a(T,F);for(let u=0;u<O.length;u+=1)O[u]&&O[u].m(F,null);a(T,Pe),a(T,H);for(let u=0;u<w.length;u+=1)w[u]&&w[u].m(H,null);D=!0},p(e,[t]){(!D||t&1)&&n!==(n=e[0].name+"")&&pe(f,n);const u={};t&8&&(u.js=`
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
    `),t&8&&(u.dart=`
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
    `),v.$set(u),(!D||t&1)&&N!==(N=e[0].name+"")&&pe(ee,N),t&6&&(Q=z(e[2]),O=Ue(O,t,Re,1,e,Q,De,F,Ne,He,null,Fe)),t&6&&(j=z(e[2]),Qe(),w=Ue(w,t,xe,1,e,j,Me,H,ze,je,null,Be),Ie())},i(e){if(!D){V(v.$$.fragment,e),V(C.$$.fragment,e),V(q.$$.fragment,e);for(let t=0;t<j.length;t+=1)V(w[t]);D=!0}},o(e){L(v.$$.fragment,e),L(C.$$.fragment,e),L(q.$$.fragment,e);for(let t=0;t<w.length;t+=1)L(w[t]);D=!1},d(e){e&&(c(l),c(d),c(b),c(_),c(A),c(P),c(Y),c(S),c(te),c(R),c(ae),c(x),c(le),c(U),c(se),c(y),c(ue),c(B),c(he),c(T)),G(v,e),G(C),G(q);for(let t=0;t<O.length;t+=1)O[t].d();for(let t=0;t<w.length;t+=1)w[t].d()}}}function tt(i,l,s){let n,{collection:f=new Ke}=l,g=200,d=[];const b=_=>s(1,g=_.code);return i.$$set=_=>{"collection"in _&&s(0,f=_.collection)},i.$$.update=()=>{i.$$.dirty&1&&s(2,d=[{code:200,body:JSON.stringify({token:"JWT_AUTH_TOKEN",record:We.dummyCollectionRecord(f),meta:{id:"abc123",name:"John Doe",username:"john.doe",email:"test@example.com",avatarUrl:"https://example.com/avatar.png",accessToken:"...",refreshToken:"...",rawUser:{}}},null,2)},{code:400,body:`
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
            `}])},s(3,n=We.getApiExampleUrl(Ge.baseUrl)),[f,g,d,n,b]}class ot extends Le{constructor(l){super(),Ee(this,l,tt,et,Je,{collection:0})}}export{ot as default};
