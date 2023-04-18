import{S as je,i as Fe,s as Ve,M as He,e as s,w as g,b as h,c as re,f as p,g as r,h as a,m as ce,x as ue,N as Pe,O as Le,k as Ee,P as Je,n as Ne,t as E,a as J,o as c,d as de,T as ze,C as Re,p as Ie,r as N,u as Ke}from"./index-7275477a.js";import{S as Qe}from"./SdkTabs-3e2121ec.js";function xe(i,l,o){const n=i.slice();return n[5]=l[o],n}function We(i,l,o){const n=i.slice();return n[5]=l[o],n}function Ue(i,l){let o,n=l[5].code+"",m,k,u,b;function _(){return l[4](l[5])}return{key:i,first:null,c(){o=s("button"),m=g(n),k=h(),p(o,"class","tab-item"),N(o,"active",l[1]===l[5].code),this.first=o},m(v,A){r(v,o,A),a(o,m),a(o,k),u||(b=Ke(o,"click",_),u=!0)},p(v,A){l=v,A&4&&n!==(n=l[5].code+"")&&ue(m,n),A&6&&N(o,"active",l[1]===l[5].code)},d(v){v&&c(o),u=!1,b()}}}function Be(i,l){let o,n,m,k;return n=new He({props:{content:l[5].body}}),{key:i,first:null,c(){o=s("div"),re(n.$$.fragment),m=h(),p(o,"class","tab-item"),N(o,"active",l[1]===l[5].code),this.first=o},m(u,b){r(u,o,b),ce(n,o,null),a(o,m),k=!0},p(u,b){l=u;const _={};b&4&&(_.content=l[5].body),n.$set(_),(!k||b&6)&&N(o,"active",l[1]===l[5].code)},i(u){k||(E(n.$$.fragment,u),k=!0)},o(u){J(n.$$.fragment,u),k=!1},d(u){u&&c(o),de(n)}}}function Ge(i){let l,o,n=i[0].name+"",m,k,u,b,_,v,A,D,z,S,j,he,F,M,pe,I,V=i[0].name+"",K,be,Q,P,G,R,X,x,Y,y,Z,fe,ee,$,te,me,ae,ke,f,ge,C,_e,ve,we,le,Oe,oe,Ae,Se,ye,se,$e,ne,W,ie,T,U,O=[],Te=new Map,Ce,B,w=[],qe=new Map,q;v=new Qe({props:{js:`
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
    `}}),C=new He({props:{content:"?expand=relField1,relField2.subRelField"}});let L=i[2];const De=e=>e[5].code;for(let e=0;e<L.length;e+=1){let t=We(i,L,e),d=De(t);Te.set(d,O[e]=Ue(d,t))}let H=i[2];const Me=e=>e[5].code;for(let e=0;e<H.length;e+=1){let t=xe(i,H,e),d=Me(t);qe.set(d,w[e]=Be(d,t))}return{c(){l=s("h3"),o=g("Auth with OAuth2 ("),m=g(n),k=g(")"),u=h(),b=s("div"),b.innerHTML=`<p>Authenticate with an OAuth2 provider and returns a new auth token and record data.</p> 
    <p>For more details please check the
        <a href="https://pocketbase.io/docs/authentication/#oauth2-integration" target="_blank" rel="noopener noreferrer">OAuth2 integration documentation
        </a>.</p>`,_=h(),re(v.$$.fragment),A=h(),D=s("h6"),D.textContent="API details",z=h(),S=s("div"),j=s("strong"),j.textContent="POST",he=h(),F=s("div"),M=s("p"),pe=g("/api/collections/"),I=s("strong"),K=g(V),be=g("/auth-with-oauth2"),Q=h(),P=s("div"),P.textContent="Body Parameters",G=h(),R=s("table"),R.innerHTML=`<thead><tr><th>Param</th> 
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
                        upload currently are not supported during OAuth2 sign-ups.</em></p></td></tr></tbody>`,X=h(),x=s("div"),x.textContent="Query parameters",Y=h(),y=s("table"),Z=s("thead"),Z.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr>`,fe=h(),ee=s("tbody"),$=s("tr"),te=s("td"),te.textContent="expand",me=h(),ae=s("td"),ae.innerHTML='<span class="label">String</span>',ke=h(),f=s("td"),ge=g(`Auto expand record relations. Ex.:
                `),re(C.$$.fragment),_e=g(`
                Supports up to 6-levels depth nested relations expansion. `),ve=s("br"),we=g(`
                The expanded relations will be appended to the record under the
                `),le=s("code"),le.textContent="expand",Oe=g(" property (eg. "),oe=s("code"),oe.textContent='"expand": {"relField1": {...}, ...}',Ae=g(`).
                `),Se=s("br"),ye=g(`
                Only the relations to which the request user has permissions to `),se=s("strong"),se.textContent="view",$e=g(" will be expanded."),ne=h(),W=s("div"),W.textContent="Responses",ie=h(),T=s("div"),U=s("div");for(let e=0;e<O.length;e+=1)O[e].c();Ce=h(),B=s("div");for(let e=0;e<w.length;e+=1)w[e].c();p(l,"class","m-b-sm"),p(b,"class","content txt-lg m-b-sm"),p(D,"class","m-b-xs"),p(j,"class","label label-primary"),p(F,"class","content"),p(S,"class","alert alert-success"),p(P,"class","section-title"),p(R,"class","table-compact table-border m-b-base"),p(x,"class","section-title"),p(y,"class","table-compact table-border m-b-base"),p(W,"class","section-title"),p(U,"class","tabs-header compact left"),p(B,"class","tabs-content"),p(T,"class","tabs")},m(e,t){r(e,l,t),a(l,o),a(l,m),a(l,k),r(e,u,t),r(e,b,t),r(e,_,t),ce(v,e,t),r(e,A,t),r(e,D,t),r(e,z,t),r(e,S,t),a(S,j),a(S,he),a(S,F),a(F,M),a(M,pe),a(M,I),a(I,K),a(M,be),r(e,Q,t),r(e,P,t),r(e,G,t),r(e,R,t),r(e,X,t),r(e,x,t),r(e,Y,t),r(e,y,t),a(y,Z),a(y,fe),a(y,ee),a(ee,$),a($,te),a($,me),a($,ae),a($,ke),a($,f),a(f,ge),ce(C,f,null),a(f,_e),a(f,ve),a(f,we),a(f,le),a(f,Oe),a(f,oe),a(f,Ae),a(f,Se),a(f,ye),a(f,se),a(f,$e),r(e,ne,t),r(e,W,t),r(e,ie,t),r(e,T,t),a(T,U);for(let d=0;d<O.length;d+=1)O[d]&&O[d].m(U,null);a(T,Ce),a(T,B);for(let d=0;d<w.length;d+=1)w[d]&&w[d].m(B,null);q=!0},p(e,[t]){(!q||t&1)&&n!==(n=e[0].name+"")&&ue(m,n);const d={};t&8&&(d.js=`
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
    `),v.$set(d),(!q||t&1)&&V!==(V=e[0].name+"")&&ue(K,V),t&6&&(L=e[2],O=Pe(O,t,De,1,e,L,Te,U,Le,Ue,null,We)),t&6&&(H=e[2],Ee(),w=Pe(w,t,Me,1,e,H,qe,B,Je,Be,null,xe),Ne())},i(e){if(!q){E(v.$$.fragment,e),E(C.$$.fragment,e);for(let t=0;t<H.length;t+=1)E(w[t]);q=!0}},o(e){J(v.$$.fragment,e),J(C.$$.fragment,e);for(let t=0;t<w.length;t+=1)J(w[t]);q=!1},d(e){e&&c(l),e&&c(u),e&&c(b),e&&c(_),de(v,e),e&&c(A),e&&c(D),e&&c(z),e&&c(S),e&&c(Q),e&&c(P),e&&c(G),e&&c(R),e&&c(X),e&&c(x),e&&c(Y),e&&c(y),de(C),e&&c(ne),e&&c(W),e&&c(ie),e&&c(T);for(let t=0;t<O.length;t+=1)O[t].d();for(let t=0;t<w.length;t+=1)w[t].d()}}}function Xe(i,l,o){let n,{collection:m=new ze}=l,k=200,u=[];const b=_=>o(1,k=_.code);return i.$$set=_=>{"collection"in _&&o(0,m=_.collection)},i.$$.update=()=>{i.$$.dirty&1&&o(2,u=[{code:200,body:JSON.stringify({token:"JWT_AUTH_TOKEN",record:Re.dummyCollectionRecord(m),meta:{id:"abc123",name:"John Doe",username:"john.doe",email:"test@example.com",avatarUrl:"https://example.com/avatar.png",accessToken:"...",refreshToken:"...",rawUser:{}}},null,2)},{code:400,body:`
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
            `}])},o(3,n=Re.getApiExampleUrl(Ie.baseUrl)),[m,k,u,n,b]}class et extends je{constructor(l){super(),Fe(this,l,Xe,Ge,Ve,{collection:0})}}export{et as default};
