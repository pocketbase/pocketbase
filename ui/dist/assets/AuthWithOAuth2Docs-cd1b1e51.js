import{S as je,i as ze,s as Fe,M as We,e as s,w,b as h,c as re,f as p,g as i,h as a,m as ce,x as de,N as De,O as Ve,k as Le,P as Ee,n as Je,t as L,a as E,o as r,d as ue,R as Ne,C as Ue,p as Ie,r as J,u as Ke}from"./index-9b05587d.js";import{S as Qe}from"./SdkTabs-088af0d8.js";function Be(c,l,o){const n=c.slice();return n[5]=l[o],n}function Me(c,l,o){const n=c.slice();return n[5]=l[o],n}function Re(c,l){let o,n=l[5].code+"",m,v,d,b;function k(){return l[4](l[5])}return{key:c,first:null,c(){o=s("button"),m=w(n),v=h(),p(o,"class","tab-item"),J(o,"active",l[1]===l[5].code),this.first=o},m(g,O){i(g,o,O),a(o,m),a(o,v),d||(b=Ke(o,"click",k),d=!0)},p(g,O){l=g,O&4&&n!==(n=l[5].code+"")&&de(m,n),O&6&&J(o,"active",l[1]===l[5].code)},d(g){g&&r(o),d=!1,b()}}}function He(c,l){let o,n,m,v;return n=new We({props:{content:l[5].body}}),{key:c,first:null,c(){o=s("div"),re(n.$$.fragment),m=h(),p(o,"class","tab-item"),J(o,"active",l[1]===l[5].code),this.first=o},m(d,b){i(d,o,b),ce(n,o,null),a(o,m),v=!0},p(d,b){l=d;const k={};b&4&&(k.content=l[5].body),n.$set(k),(!v||b&6)&&J(o,"active",l[1]===l[5].code)},i(d){v||(L(n.$$.fragment,d),v=!0)},o(d){E(n.$$.fragment,d),v=!1},d(d){d&&r(o),ue(n)}}}function Ge(c){let l,o,n=c[0].name+"",m,v,d,b,k,g,O,x,N,A,j,he,z,P,pe,I,F=c[0].name+"",K,be,Q,D,G,U,X,B,Y,S,Z,fe,ee,T,te,me,ae,ve,f,we,C,ke,ge,_e,le,ye,oe,Oe,Ae,Se,se,Te,ne,M,ie,$,R,y=[],$e=new Map,Ce,H,_=[],qe=new Map,q;g=new Qe({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${c[3]}');

        ...

        // This method initializes a one-off realtime subscription and will
        // open a popup window with the OAuth2 vendor page to authenticate.
        //
        // Once the external OAuth2 sign-in/sign-up flow is completed, the popup
        // window will be automatically closed and the OAuth2 data sent back
        // to the user through the previously established realtime connection.
        const authData = await pb.collection('users').authWithOAuth2({ provider: 'google' });

        // after the above you can also access the auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.model.id);

        // "logout" the last authenticated model
        pb.authStore.clear();
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';
        import 'package:url_launcher/url_launcher.dart';

        final pb = PocketBase('${c[3]}');

        ...

        // This method initializes a one-off realtime subscription and will
        // call the provided urlCallback with the OAuth2 vendor url to authenticate.
        //
        // Once the external OAuth2 sign-in/sign-up flow is completed, the browser
        // window will be automatically closed and the OAuth2 data sent back
        // to the user through the previously established realtime connection.
        final authData = await pb.collection('users').authWithOAuth2('google', (url) async {
          await launchUrl(url);
        });

        // after the above you can also access the auth data from the authStore
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.model.id);

        // "logout" the last authenticated model
        pb.authStore.clear();
    `}}),C=new We({props:{content:"?expand=relField1,relField2.subRelField"}});let V=c[2];const xe=e=>e[5].code;for(let e=0;e<V.length;e+=1){let t=Me(c,V,e),u=xe(t);$e.set(u,y[e]=Re(u,t))}let W=c[2];const Pe=e=>e[5].code;for(let e=0;e<W.length;e+=1){let t=Be(c,W,e),u=Pe(t);qe.set(u,_[e]=He(u,t))}return{c(){l=s("h3"),o=w("Auth with OAuth2 ("),m=w(n),v=w(")"),d=h(),b=s("div"),b.innerHTML=`<p>Authenticate with an OAuth2 provider and returns a new auth token and record data.</p> 
    <p>For more details please check the
        <a href="https://pocketbase.io/docs/authentication/#oauth2-integration" target="_blank" rel="noopener noreferrer">OAuth2 integration documentation
        </a>.</p>`,k=h(),re(g.$$.fragment),O=h(),x=s("h6"),x.textContent="API details",N=h(),A=s("div"),j=s("strong"),j.textContent="POST",he=h(),z=s("div"),P=s("p"),pe=w("/api/collections/"),I=s("strong"),K=w(F),be=w("/auth-with-oauth2"),Q=h(),D=s("div"),D.textContent="Body Parameters",G=h(),U=s("table"),U.innerHTML=`<thead><tr><th>Param</th> 
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
                        upload currently are not supported during OAuth2 sign-ups.</em></p></td></tr></tbody>`,X=h(),B=s("div"),B.textContent="Query parameters",Y=h(),S=s("table"),Z=s("thead"),Z.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr>`,fe=h(),ee=s("tbody"),T=s("tr"),te=s("td"),te.textContent="expand",me=h(),ae=s("td"),ae.innerHTML='<span class="label">String</span>',ve=h(),f=s("td"),we=w(`Auto expand record relations. Ex.:
                `),re(C.$$.fragment),ke=w(`
                Supports up to 6-levels depth nested relations expansion. `),ge=s("br"),_e=w(`
                The expanded relations will be appended to the record under the
                `),le=s("code"),le.textContent="expand",ye=w(" property (eg. "),oe=s("code"),oe.textContent='"expand": {"relField1": {...}, ...}',Oe=w(`).
                `),Ae=s("br"),Se=w(`
                Only the relations to which the request user has permissions to `),se=s("strong"),se.textContent="view",Te=w(" will be expanded."),ne=h(),M=s("div"),M.textContent="Responses",ie=h(),$=s("div"),R=s("div");for(let e=0;e<y.length;e+=1)y[e].c();Ce=h(),H=s("div");for(let e=0;e<_.length;e+=1)_[e].c();p(l,"class","m-b-sm"),p(b,"class","content txt-lg m-b-sm"),p(x,"class","m-b-xs"),p(j,"class","label label-primary"),p(z,"class","content"),p(A,"class","alert alert-success"),p(D,"class","section-title"),p(U,"class","table-compact table-border m-b-base"),p(B,"class","section-title"),p(S,"class","table-compact table-border m-b-base"),p(M,"class","section-title"),p(R,"class","tabs-header compact left"),p(H,"class","tabs-content"),p($,"class","tabs")},m(e,t){i(e,l,t),a(l,o),a(l,m),a(l,v),i(e,d,t),i(e,b,t),i(e,k,t),ce(g,e,t),i(e,O,t),i(e,x,t),i(e,N,t),i(e,A,t),a(A,j),a(A,he),a(A,z),a(z,P),a(P,pe),a(P,I),a(I,K),a(P,be),i(e,Q,t),i(e,D,t),i(e,G,t),i(e,U,t),i(e,X,t),i(e,B,t),i(e,Y,t),i(e,S,t),a(S,Z),a(S,fe),a(S,ee),a(ee,T),a(T,te),a(T,me),a(T,ae),a(T,ve),a(T,f),a(f,we),ce(C,f,null),a(f,ke),a(f,ge),a(f,_e),a(f,le),a(f,ye),a(f,oe),a(f,Oe),a(f,Ae),a(f,Se),a(f,se),a(f,Te),i(e,ne,t),i(e,M,t),i(e,ie,t),i(e,$,t),a($,R);for(let u=0;u<y.length;u+=1)y[u]&&y[u].m(R,null);a($,Ce),a($,H);for(let u=0;u<_.length;u+=1)_[u]&&_[u].m(H,null);q=!0},p(e,[t]){(!q||t&1)&&n!==(n=e[0].name+"")&&de(m,n);const u={};t&8&&(u.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        // This method initializes a one-off realtime subscription and will
        // open a popup window with the OAuth2 vendor page to authenticate.
        //
        // Once the external OAuth2 sign-in/sign-up flow is completed, the popup
        // window will be automatically closed and the OAuth2 data sent back
        // to the user through the previously established realtime connection.
        const authData = await pb.collection('users').authWithOAuth2({ provider: 'google' });

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

        // This method initializes a one-off realtime subscription and will
        // call the provided urlCallback with the OAuth2 vendor url to authenticate.
        //
        // Once the external OAuth2 sign-in/sign-up flow is completed, the browser
        // window will be automatically closed and the OAuth2 data sent back
        // to the user through the previously established realtime connection.
        final authData = await pb.collection('users').authWithOAuth2('google', (url) async {
          await launchUrl(url);
        });

        // after the above you can also access the auth data from the authStore
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.model.id);

        // "logout" the last authenticated model
        pb.authStore.clear();
    `),g.$set(u),(!q||t&1)&&F!==(F=e[0].name+"")&&de(K,F),t&6&&(V=e[2],y=De(y,t,xe,1,e,V,$e,R,Ve,Re,null,Me)),t&6&&(W=e[2],Le(),_=De(_,t,Pe,1,e,W,qe,H,Ee,He,null,Be),Je())},i(e){if(!q){L(g.$$.fragment,e),L(C.$$.fragment,e);for(let t=0;t<W.length;t+=1)L(_[t]);q=!0}},o(e){E(g.$$.fragment,e),E(C.$$.fragment,e);for(let t=0;t<_.length;t+=1)E(_[t]);q=!1},d(e){e&&r(l),e&&r(d),e&&r(b),e&&r(k),ue(g,e),e&&r(O),e&&r(x),e&&r(N),e&&r(A),e&&r(Q),e&&r(D),e&&r(G),e&&r(U),e&&r(X),e&&r(B),e&&r(Y),e&&r(S),ue(C),e&&r(ne),e&&r(M),e&&r(ie),e&&r($);for(let t=0;t<y.length;t+=1)y[t].d();for(let t=0;t<_.length;t+=1)_[t].d()}}}function Xe(c,l,o){let n,{collection:m=new Ne}=l,v=200,d=[];const b=k=>o(1,v=k.code);return c.$$set=k=>{"collection"in k&&o(0,m=k.collection)},c.$$.update=()=>{c.$$.dirty&1&&o(2,d=[{code:200,body:JSON.stringify({token:"JWT_AUTH_TOKEN",record:Ue.dummyCollectionRecord(m),meta:{id:"abc123",name:"John Doe",username:"john.doe",email:"test@example.com",avatarUrl:"https://example.com/avatar.png",accessToken:"...",refreshToken:"...",rawUser:{}}},null,2)},{code:400,body:`
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
            `}])},o(3,n=Ue.getApiExampleUrl(Ie.baseUrl)),[m,v,d,n,b]}class et extends je{constructor(l){super(),ze(this,l,Xe,Ge,Fe,{collection:0})}}export{et as default};
