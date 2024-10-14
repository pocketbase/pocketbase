import{S as xe,i as Ee,s as Je,Q as Qe,R as je,T as z,e as o,w as k,b as h,c as I,f as p,g as r,h as a,m as K,x as pe,U as Ue,V as Ne,k as ze,W as Ie,n as Ke,t as j,a as x,o as c,d as G,C as Be,p as Ge,r as X,u as Xe}from"./index-BxtrQeCj.js";import{F as Ye}from"./FieldsQueryParam-Dpe8NKnz.js";function Fe(s,l,n){const i=s.slice();return i[5]=l[n],i}function Le(s,l,n){const i=s.slice();return i[5]=l[n],i}function He(s,l){let n,i=l[5].code+"",f,g,d,b;function _(){return l[4](l[5])}return{key:s,first:null,c(){n=o("button"),f=k(i),g=h(),p(n,"class","tab-item"),X(n,"active",l[1]===l[5].code),this.first=n},m(v,O){r(v,n,O),a(n,f),a(n,g),d||(b=Xe(n,"click",_),d=!0)},p(v,O){l=v,O&4&&i!==(i=l[5].code+"")&&pe(f,i),O&6&&X(n,"active",l[1]===l[5].code)},d(v){v&&c(n),d=!1,b()}}}function Ve(s,l){let n,i,f,g;return i=new je({props:{content:l[5].body}}),{key:s,first:null,c(){n=o("div"),I(i.$$.fragment),f=h(),p(n,"class","tab-item"),X(n,"active",l[1]===l[5].code),this.first=n},m(d,b){r(d,n,b),K(i,n,null),a(n,f),g=!0},p(d,b){l=d;const _={};b&4&&(_.content=l[5].body),i.$set(_),(!g||b&6)&&X(n,"active",l[1]===l[5].code)},i(d){g||(j(i.$$.fragment,d),g=!0)},o(d){x(i.$$.fragment,d),g=!1},d(d){d&&c(n),G(i)}}}function Ze(s){let l,n,i=s[0].name+"",f,g,d,b,_,v,O,D,Y,A,E,be,J,P,me,Z,Q=s[0].name+"",ee,fe,te,M,ae,W,le,U,ne,y,oe,ge,B,S,se,ke,ie,_e,m,ve,C,we,$e,Oe,re,Ae,ce,ye,Se,Te,de,Ce,qe,q,ue,F,he,T,L,$=[],Re=new Map,De,H,w=[],Pe=new Map,R;v=new Qe({props:{js:`
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
    `}}),C=new je({props:{content:"?expand=relField1,relField2.subRelField"}}),q=new Ye({props:{prefix:"record."}});let N=z(s[2]);const Me=e=>e[5].code;for(let e=0;e<N.length;e+=1){let t=Le(s,N,e),u=Me(t);Re.set(u,$[e]=He(u,t))}let V=z(s[2]);const We=e=>e[5].code;for(let e=0;e<V.length;e+=1){let t=Fe(s,V,e),u=We(t);Pe.set(u,w[e]=Ve(u,t))}return{c(){l=o("h3"),n=k("Auth with OAuth2 ("),f=k(i),g=k(")"),d=h(),b=o("div"),b.innerHTML=`<p>Authenticate with an OAuth2 provider and returns a new auth token and record data.</p> <p>For more details please check the
        <a href="https://pocketbase.io/docs/authentication/#oauth2-integration" target="_blank" rel="noopener noreferrer">OAuth2 integration documentation
        </a>.</p>`,_=h(),I(v.$$.fragment),O=h(),D=o("h6"),D.textContent="API details",Y=h(),A=o("div"),E=o("strong"),E.textContent="POST",be=h(),J=o("div"),P=o("p"),me=k("/api/collections/"),Z=o("strong"),ee=k(Q),fe=k("/auth-with-oauth2"),te=h(),M=o("div"),M.textContent="Body Parameters",ae=h(),W=o("table"),W.innerHTML=`<thead><tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr></thead> <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>provider</span></div></td> <td><span class="label">String</span></td> <td>The name of the OAuth2 client provider (eg. &quot;google&quot;).</td></tr> <tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>code</span></div></td> <td><span class="label">String</span></td> <td>The authorization code returned from the initial request.</td></tr> <tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>codeVerifier</span></div></td> <td><span class="label">String</span></td> <td>The code verifier sent with the initial request as part of the code_challenge.</td></tr> <tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>redirectURL</span></div></td> <td><span class="label">String</span></td> <td>The redirect url sent with the initial request.</td></tr> <tr><td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>createData</span></div></td> <td><span class="label">Object</span></td> <td><p>Optional data that will be used when creating the auth record on OAuth2 sign-up.</p> <p>The created auth record must comply with the same requirements and validations in the
                    regular <strong>create</strong> action.
                    <br/> <em>The data can only be in <code>json</code>, aka. <code>multipart/form-data</code> and files
                        upload currently are not supported during OAuth2 sign-ups.</em></p></td></tr></tbody>`,le=h(),U=o("div"),U.textContent="Query parameters",ne=h(),y=o("table"),oe=o("thead"),oe.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',ge=h(),B=o("tbody"),S=o("tr"),se=o("td"),se.textContent="expand",ke=h(),ie=o("td"),ie.innerHTML='<span class="label">String</span>',_e=h(),m=o("td"),ve=k(`Auto expand record relations. Ex.:
                `),I(C.$$.fragment),we=k(`
                Supports up to 6-levels depth nested relations expansion. `),$e=o("br"),Oe=k(`
                The expanded relations will be appended to the record under the
                `),re=o("code"),re.textContent="expand",Ae=k(" property (eg. "),ce=o("code"),ce.textContent='"expand": {"relField1": {...}, ...}',ye=k(`).
                `),Se=o("br"),Te=k(`
                Only the relations to which the request user has permissions to `),de=o("strong"),de.textContent="view",Ce=k(" will be expanded."),qe=h(),I(q.$$.fragment),ue=h(),F=o("div"),F.textContent="Responses",he=h(),T=o("div"),L=o("div");for(let e=0;e<$.length;e+=1)$[e].c();De=h(),H=o("div");for(let e=0;e<w.length;e+=1)w[e].c();p(l,"class","m-b-sm"),p(b,"class","content txt-lg m-b-sm"),p(D,"class","m-b-xs"),p(E,"class","label label-primary"),p(J,"class","content"),p(A,"class","alert alert-success"),p(M,"class","section-title"),p(W,"class","table-compact table-border m-b-base"),p(U,"class","section-title"),p(y,"class","table-compact table-border m-b-base"),p(F,"class","section-title"),p(L,"class","tabs-header compact combined left"),p(H,"class","tabs-content"),p(T,"class","tabs")},m(e,t){r(e,l,t),a(l,n),a(l,f),a(l,g),r(e,d,t),r(e,b,t),r(e,_,t),K(v,e,t),r(e,O,t),r(e,D,t),r(e,Y,t),r(e,A,t),a(A,E),a(A,be),a(A,J),a(J,P),a(P,me),a(P,Z),a(Z,ee),a(P,fe),r(e,te,t),r(e,M,t),r(e,ae,t),r(e,W,t),r(e,le,t),r(e,U,t),r(e,ne,t),r(e,y,t),a(y,oe),a(y,ge),a(y,B),a(B,S),a(S,se),a(S,ke),a(S,ie),a(S,_e),a(S,m),a(m,ve),K(C,m,null),a(m,we),a(m,$e),a(m,Oe),a(m,re),a(m,Ae),a(m,ce),a(m,ye),a(m,Se),a(m,Te),a(m,de),a(m,Ce),a(B,qe),K(q,B,null),r(e,ue,t),r(e,F,t),r(e,he,t),r(e,T,t),a(T,L);for(let u=0;u<$.length;u+=1)$[u]&&$[u].m(L,null);a(T,De),a(T,H);for(let u=0;u<w.length;u+=1)w[u]&&w[u].m(H,null);R=!0},p(e,[t]){(!R||t&1)&&i!==(i=e[0].name+"")&&pe(f,i);const u={};t&9&&(u.js=`
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
    `),v.$set(u),(!R||t&1)&&Q!==(Q=e[0].name+"")&&pe(ee,Q),t&6&&(N=z(e[2]),$=Ue($,t,Me,1,e,N,Re,L,Ne,He,null,Le)),t&6&&(V=z(e[2]),ze(),w=Ue(w,t,We,1,e,V,Pe,H,Ie,Ve,null,Fe),Ke())},i(e){if(!R){j(v.$$.fragment,e),j(C.$$.fragment,e),j(q.$$.fragment,e);for(let t=0;t<V.length;t+=1)j(w[t]);R=!0}},o(e){x(v.$$.fragment,e),x(C.$$.fragment,e),x(q.$$.fragment,e);for(let t=0;t<w.length;t+=1)x(w[t]);R=!1},d(e){e&&(c(l),c(d),c(b),c(_),c(O),c(D),c(Y),c(A),c(te),c(M),c(ae),c(W),c(le),c(U),c(ne),c(y),c(ue),c(F),c(he),c(T)),G(v,e),G(C),G(q);for(let t=0;t<$.length;t+=1)$[t].d();for(let t=0;t<w.length;t+=1)w[t].d()}}}function et(s,l,n){let i,{collection:f}=l,g=200,d=[];const b=_=>n(1,g=_.code);return s.$$set=_=>{"collection"in _&&n(0,f=_.collection)},s.$$.update=()=>{s.$$.dirty&1&&n(2,d=[{code:200,body:JSON.stringify({token:"JWT_AUTH_TOKEN",record:Be.dummyCollectionRecord(f),meta:{id:"abc123",name:"John Doe",username:"john.doe",email:"test@example.com",avatarURL:"https://example.com/avatar.png",accessToken:"...",refreshToken:"...",rawUser:{}}},null,2)},{code:400,body:`
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
            `}])},n(3,i=Be.getApiExampleUrl(Ge.baseURL)),[f,g,d,i,b]}class lt extends xe{constructor(l){super(),Ee(this,l,et,Ze,Je,{collection:0})}}export{lt as default};
