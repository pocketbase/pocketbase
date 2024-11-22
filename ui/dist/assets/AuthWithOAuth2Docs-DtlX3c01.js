import{S as Ee,i as Je,s as xe,U as Ie,V as Ve,W as Q,f as o,y as k,h,c as z,j as p,l as r,n as a,m as G,G as pe,X as Ue,Y as Ne,D as Qe,Z as ze,E as Ge,t as V,a as E,u as c,d as K,I as Be,p as Ke,k as X,o as Xe}from"./index-BASVPtys.js";import{F as Ye}from"./FieldsQueryParam-CHh5pfbg.js";function Fe(s,l,n){const i=s.slice();return i[5]=l[n],i}function Le(s,l,n){const i=s.slice();return i[5]=l[n],i}function je(s,l){let n,i=l[5].code+"",f,g,d,b;function _(){return l[4](l[5])}return{key:s,first:null,c(){n=o("button"),f=k(i),g=h(),p(n,"class","tab-item"),X(n,"active",l[1]===l[5].code),this.first=n},m(v,O){r(v,n,O),a(n,f),a(n,g),d||(b=Xe(n,"click",_),d=!0)},p(v,O){l=v,O&4&&i!==(i=l[5].code+"")&&pe(f,i),O&6&&X(n,"active",l[1]===l[5].code)},d(v){v&&c(n),d=!1,b()}}}function He(s,l){let n,i,f,g;return i=new Ve({props:{content:l[5].body}}),{key:s,first:null,c(){n=o("div"),z(i.$$.fragment),f=h(),p(n,"class","tab-item"),X(n,"active",l[1]===l[5].code),this.first=n},m(d,b){r(d,n,b),G(i,n,null),a(n,f),g=!0},p(d,b){l=d;const _={};b&4&&(_.content=l[5].body),i.$set(_),(!g||b&6)&&X(n,"active",l[1]===l[5].code)},i(d){g||(V(i.$$.fragment,d),g=!0)},o(d){E(i.$$.fragment,d),g=!1},d(d){d&&c(n),K(i)}}}function Ze(s){let l,n,i=s[0].name+"",f,g,d,b,_,v,O,R,Y,y,J,be,x,P,me,Z,I=s[0].name+"",ee,fe,te,M,ae,W,le,U,ne,A,oe,ge,B,S,se,ke,ie,_e,m,ve,C,we,$e,Oe,re,ye,ce,Ae,Se,Te,de,Ce,qe,q,ue,F,he,T,L,$=[],De=new Map,Re,j,w=[],Pe=new Map,D;v=new Ie({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${s[3]}');

        ...

        // OAuth2 authentication with a single realtime call.
        //
        // Make sure to register ${s[3]}/api/oauth2-redirect as redirect url.
        const authData = await pb.collection('${s[0].name}').authWithOAuth2({ provider: 'google' });

        // OR authenticate with manual OAuth2 code exchange
        // const authData = await pb.collection('${s[0].name}').authWithOAuth2Code(...);

        // after the above you can also access the auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.record.id);

        // "logout"
        pb.authStore.clear();
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';
        import 'package:url_launcher/url_launcher.dart';

        final pb = PocketBase('${s[3]}');

        ...

        // OAuth2 authentication with a single realtime call.
        //
        // Make sure to register ${s[3]}/api/oauth2-redirect as redirect url.
        final authData = await pb.collection('${s[0].name}').authWithOAuth2('google', (url) async {
          await launchUrl(url);
        });

        // OR authenticate with manual OAuth2 code exchange
        // final authData = await pb.collection('${s[0].name}').authWithOAuth2Code(...);

        // after the above you can also access the auth data from the authStore
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.record.id);

        // "logout"
        pb.authStore.clear();
    `}}),C=new Ve({props:{content:"?expand=relField1,relField2.subRelField"}}),q=new Ye({props:{prefix:"record."}});let N=Q(s[2]);const Me=e=>e[5].code;for(let e=0;e<N.length;e+=1){let t=Le(s,N,e),u=Me(t);De.set(u,$[e]=je(u,t))}let H=Q(s[2]);const We=e=>e[5].code;for(let e=0;e<H.length;e+=1){let t=Fe(s,H,e),u=We(t);Pe.set(u,w[e]=He(u,t))}return{c(){l=o("h3"),n=k("Auth with OAuth2 ("),f=k(i),g=k(")"),d=h(),b=o("div"),b.innerHTML=`<p>Authenticate with an OAuth2 provider and returns a new auth token and record data.</p> <p>For more details please check the
        <a href="https://pocketbase.io/docs/authentication/#authenticate-with-oauth2" target="_blank" rel="noopener noreferrer">OAuth2 integration documentation
        </a>.</p>`,_=h(),z(v.$$.fragment),O=h(),R=o("h6"),R.textContent="API details",Y=h(),y=o("div"),J=o("strong"),J.textContent="POST",be=h(),x=o("div"),P=o("p"),me=k("/api/collections/"),Z=o("strong"),ee=k(I),fe=k("/auth-with-oauth2"),te=h(),M=o("div"),M.textContent="Body Parameters",ae=h(),W=o("table"),W.innerHTML=`<thead><tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr></thead> <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>provider</span></div></td> <td><span class="label">String</span></td> <td>The name of the OAuth2 client provider (eg. &quot;google&quot;).</td></tr> <tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>code</span></div></td> <td><span class="label">String</span></td> <td>The authorization code returned from the initial request.</td></tr> <tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>codeVerifier</span></div></td> <td><span class="label">String</span></td> <td>The code verifier sent with the initial request as part of the code_challenge.</td></tr> <tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>redirectURL</span></div></td> <td><span class="label">String</span></td> <td>The redirect url sent with the initial request.</td></tr> <tr><td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>createData</span></div></td> <td><span class="label">Object</span></td> <td><p>Optional data that will be used when creating the auth record on OAuth2 sign-up.</p> <p>The created auth record must comply with the same requirements and validations in the
                    regular <strong>create</strong> action.
                    <br/> <em>The data can only be in <code>json</code>, aka. <code>multipart/form-data</code> and files
                        upload currently are not supported during OAuth2 sign-ups.</em></p></td></tr></tbody>`,le=h(),U=o("div"),U.textContent="Query parameters",ne=h(),A=o("table"),oe=o("thead"),oe.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',ge=h(),B=o("tbody"),S=o("tr"),se=o("td"),se.textContent="expand",ke=h(),ie=o("td"),ie.innerHTML='<span class="label">String</span>',_e=h(),m=o("td"),ve=k(`Auto expand record relations. Ex.:
                `),z(C.$$.fragment),we=k(`
                Supports up to 6-levels depth nested relations expansion. `),$e=o("br"),Oe=k(`
                The expanded relations will be appended to the record under the
                `),re=o("code"),re.textContent="expand",ye=k(" property (eg. "),ce=o("code"),ce.textContent='"expand": {"relField1": {...}, ...}',Ae=k(`).
                `),Se=o("br"),Te=k(`
                Only the relations to which the request user has permissions to `),de=o("strong"),de.textContent="view",Ce=k(" will be expanded."),qe=h(),z(q.$$.fragment),ue=h(),F=o("div"),F.textContent="Responses",he=h(),T=o("div"),L=o("div");for(let e=0;e<$.length;e+=1)$[e].c();Re=h(),j=o("div");for(let e=0;e<w.length;e+=1)w[e].c();p(l,"class","m-b-sm"),p(b,"class","content txt-lg m-b-sm"),p(R,"class","m-b-xs"),p(J,"class","label label-primary"),p(x,"class","content"),p(y,"class","alert alert-success"),p(M,"class","section-title"),p(W,"class","table-compact table-border m-b-base"),p(U,"class","section-title"),p(A,"class","table-compact table-border m-b-base"),p(F,"class","section-title"),p(L,"class","tabs-header compact combined left"),p(j,"class","tabs-content"),p(T,"class","tabs")},m(e,t){r(e,l,t),a(l,n),a(l,f),a(l,g),r(e,d,t),r(e,b,t),r(e,_,t),G(v,e,t),r(e,O,t),r(e,R,t),r(e,Y,t),r(e,y,t),a(y,J),a(y,be),a(y,x),a(x,P),a(P,me),a(P,Z),a(Z,ee),a(P,fe),r(e,te,t),r(e,M,t),r(e,ae,t),r(e,W,t),r(e,le,t),r(e,U,t),r(e,ne,t),r(e,A,t),a(A,oe),a(A,ge),a(A,B),a(B,S),a(S,se),a(S,ke),a(S,ie),a(S,_e),a(S,m),a(m,ve),G(C,m,null),a(m,we),a(m,$e),a(m,Oe),a(m,re),a(m,ye),a(m,ce),a(m,Ae),a(m,Se),a(m,Te),a(m,de),a(m,Ce),a(B,qe),G(q,B,null),r(e,ue,t),r(e,F,t),r(e,he,t),r(e,T,t),a(T,L);for(let u=0;u<$.length;u+=1)$[u]&&$[u].m(L,null);a(T,Re),a(T,j);for(let u=0;u<w.length;u+=1)w[u]&&w[u].m(j,null);D=!0},p(e,[t]){(!D||t&1)&&i!==(i=e[0].name+"")&&pe(f,i);const u={};t&9&&(u.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        // OAuth2 authentication with a single realtime call.
        //
        // Make sure to register ${e[3]}/api/oauth2-redirect as redirect url.
        const authData = await pb.collection('${e[0].name}').authWithOAuth2({ provider: 'google' });

        // OR authenticate with manual OAuth2 code exchange
        // const authData = await pb.collection('${e[0].name}').authWithOAuth2Code(...);

        // after the above you can also access the auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.record.id);

        // "logout"
        pb.authStore.clear();
    `),t&9&&(u.dart=`
        import 'package:pocketbase/pocketbase.dart';
        import 'package:url_launcher/url_launcher.dart';

        final pb = PocketBase('${e[3]}');

        ...

        // OAuth2 authentication with a single realtime call.
        //
        // Make sure to register ${e[3]}/api/oauth2-redirect as redirect url.
        final authData = await pb.collection('${e[0].name}').authWithOAuth2('google', (url) async {
          await launchUrl(url);
        });

        // OR authenticate with manual OAuth2 code exchange
        // final authData = await pb.collection('${e[0].name}').authWithOAuth2Code(...);

        // after the above you can also access the auth data from the authStore
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.record.id);

        // "logout"
        pb.authStore.clear();
    `),v.$set(u),(!D||t&1)&&I!==(I=e[0].name+"")&&pe(ee,I),t&6&&(N=Q(e[2]),$=Ue($,t,Me,1,e,N,De,L,Ne,je,null,Le)),t&6&&(H=Q(e[2]),Qe(),w=Ue(w,t,We,1,e,H,Pe,j,ze,He,null,Fe),Ge())},i(e){if(!D){V(v.$$.fragment,e),V(C.$$.fragment,e),V(q.$$.fragment,e);for(let t=0;t<H.length;t+=1)V(w[t]);D=!0}},o(e){E(v.$$.fragment,e),E(C.$$.fragment,e),E(q.$$.fragment,e);for(let t=0;t<w.length;t+=1)E(w[t]);D=!1},d(e){e&&(c(l),c(d),c(b),c(_),c(O),c(R),c(Y),c(y),c(te),c(M),c(ae),c(W),c(le),c(U),c(ne),c(A),c(ue),c(F),c(he),c(T)),K(v,e),K(C),K(q);for(let t=0;t<$.length;t+=1)$[t].d();for(let t=0;t<w.length;t+=1)w[t].d()}}}function et(s,l,n){let i,{collection:f}=l,g=200,d=[];const b=_=>n(1,g=_.code);return s.$$set=_=>{"collection"in _&&n(0,f=_.collection)},s.$$.update=()=>{s.$$.dirty&1&&n(2,d=[{code:200,body:JSON.stringify({token:"JWT_AUTH_TOKEN",record:Be.dummyCollectionRecord(f),meta:{id:"abc123",name:"John Doe",username:"john.doe",email:"test@example.com",avatarURL:"https://example.com/avatar.png",accessToken:"...",refreshToken:"...",rawUser:{}}},null,2)},{code:400,body:`
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
            `}])},n(3,i=Be.getApiExampleUrl(Ke.baseURL)),[f,g,d,i,b]}class lt extends Ee{constructor(l){super(),Je(this,l,et,Ze,xe,{collection:0})}}export{lt as default};
