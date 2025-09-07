import{S as Je,i as xe,s as Ee,V as Ne,W as je,X as Q,h as r,d as Z,t as j,a as J,I as pe,Z as Ue,_ as Ie,C as Qe,$ as Ze,D as ze,l as c,n as a,m as z,u as o,A as _,v as h,c as K,w as p,J as Be,p as Ke,k as X,o as Xe}from"./index-CkwOC79g.js";import{F as Ge}from"./FieldsQueryParam-haBHHZQ1.js";function Fe(s,l,n){const i=s.slice();return i[5]=l[n],i}function Le(s,l,n){const i=s.slice();return i[5]=l[n],i}function He(s,l){let n,i=l[5].code+"",f,g,d,b;function k(){return l[4](l[5])}return{key:s,first:null,c(){n=o("button"),f=_(i),g=h(),p(n,"class","tab-item"),X(n,"active",l[1]===l[5].code),this.first=n},m(v,O){c(v,n,O),a(n,f),a(n,g),d||(b=Xe(n,"click",k),d=!0)},p(v,O){l=v,O&4&&i!==(i=l[5].code+"")&&pe(f,i),O&6&&X(n,"active",l[1]===l[5].code)},d(v){v&&r(n),d=!1,b()}}}function Ve(s,l){let n,i,f,g;return i=new je({props:{content:l[5].body}}),{key:s,first:null,c(){n=o("div"),K(i.$$.fragment),f=h(),p(n,"class","tab-item"),X(n,"active",l[1]===l[5].code),this.first=n},m(d,b){c(d,n,b),z(i,n,null),a(n,f),g=!0},p(d,b){l=d;const k={};b&4&&(k.content=l[5].body),i.$set(k),(!g||b&6)&&X(n,"active",l[1]===l[5].code)},i(d){g||(J(i.$$.fragment,d),g=!0)},o(d){j(i.$$.fragment,d),g=!1},d(d){d&&r(n),Z(i)}}}function Ye(s){let l,n,i=s[0].name+"",f,g,d,b,k,v,O,R,G,A,x,be,E,P,me,Y,N=s[0].name+"",ee,fe,te,M,ae,W,le,U,ne,y,oe,ge,B,S,se,_e,ie,ke,m,ve,C,we,$e,Oe,re,Ae,ce,ye,Se,Te,de,Ce,qe,q,ue,F,he,T,L,$=[],De=new Map,Re,H,w=[],Pe=new Map,D;v=new Ne({props:{js:`
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
    `}}),C=new je({props:{content:"?expand=relField1,relField2.subRelField"}}),q=new Ge({props:{prefix:"record."}});let I=Q(s[2]);const Me=e=>e[5].code;for(let e=0;e<I.length;e+=1){let t=Le(s,I,e),u=Me(t);De.set(u,$[e]=He(u,t))}let V=Q(s[2]);const We=e=>e[5].code;for(let e=0;e<V.length;e+=1){let t=Fe(s,V,e),u=We(t);Pe.set(u,w[e]=Ve(u,t))}return{c(){l=o("h3"),n=_("Auth with OAuth2 ("),f=_(i),g=_(")"),d=h(),b=o("div"),b.innerHTML=`<p>Authenticate with an OAuth2 provider and returns a new auth token and record data.</p> <p>For more details please check the
        <a href="https://pocketbase.io/docs/authentication#authenticate-with-oauth2" target="_blank" rel="noopener noreferrer">OAuth2 integration documentation
        </a>.</p>`,k=h(),K(v.$$.fragment),O=h(),R=o("h6"),R.textContent="API details",G=h(),A=o("div"),x=o("strong"),x.textContent="POST",be=h(),E=o("div"),P=o("p"),me=_("/api/collections/"),Y=o("strong"),ee=_(N),fe=_("/auth-with-oauth2"),te=h(),M=o("div"),M.textContent="Body Parameters",ae=h(),W=o("table"),W.innerHTML=`<thead><tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr></thead> <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>provider</span></div></td> <td><span class="label">String</span></td> <td>The name of the OAuth2 client provider (eg. &quot;google&quot;).</td></tr> <tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>code</span></div></td> <td><span class="label">String</span></td> <td>The authorization code returned from the initial request.</td></tr> <tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>codeVerifier</span></div></td> <td><span class="label">String</span></td> <td>The code verifier sent with the initial request as part of the code_challenge.</td></tr> <tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>redirectURL</span></div></td> <td><span class="label">String</span></td> <td>The redirect url sent with the initial request.</td></tr> <tr><td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>createData</span></div></td> <td><span class="label">Object</span></td> <td><p>Optional data that will be used when creating the auth record on OAuth2 sign-up.</p> <p>The created auth record must comply with the same requirements and validations in the
                    regular <strong>create</strong> action.
                    <br/> <em>The data can only be in <code>json</code>, aka. <code>multipart/form-data</code> and files
                        upload currently are not supported during OAuth2 sign-ups.</em></p></td></tr></tbody>`,le=h(),U=o("div"),U.textContent="Query parameters",ne=h(),y=o("table"),oe=o("thead"),oe.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',ge=h(),B=o("tbody"),S=o("tr"),se=o("td"),se.textContent="expand",_e=h(),ie=o("td"),ie.innerHTML='<span class="label">String</span>',ke=h(),m=o("td"),ve=_(`Auto expand record relations. Ex.:
                `),K(C.$$.fragment),we=_(`
                Supports up to 6-levels depth nested relations expansion. `),$e=o("br"),Oe=_(`
                The expanded relations will be appended to the record under the
                `),re=o("code"),re.textContent="expand",Ae=_(" property (eg. "),ce=o("code"),ce.textContent='"expand": {"relField1": {...}, ...}',ye=_(`).
                `),Se=o("br"),Te=_(`
                Only the relations to which the request user has permissions to `),de=o("strong"),de.textContent="view",Ce=_(" will be expanded."),qe=h(),K(q.$$.fragment),ue=h(),F=o("div"),F.textContent="Responses",he=h(),T=o("div"),L=o("div");for(let e=0;e<$.length;e+=1)$[e].c();Re=h(),H=o("div");for(let e=0;e<w.length;e+=1)w[e].c();p(l,"class","m-b-sm"),p(b,"class","content txt-lg m-b-sm"),p(R,"class","m-b-xs"),p(x,"class","label label-primary"),p(E,"class","content"),p(A,"class","alert alert-success"),p(M,"class","section-title"),p(W,"class","table-compact table-border m-b-base"),p(U,"class","section-title"),p(y,"class","table-compact table-border m-b-base"),p(F,"class","section-title"),p(L,"class","tabs-header compact combined left"),p(H,"class","tabs-content"),p(T,"class","tabs")},m(e,t){c(e,l,t),a(l,n),a(l,f),a(l,g),c(e,d,t),c(e,b,t),c(e,k,t),z(v,e,t),c(e,O,t),c(e,R,t),c(e,G,t),c(e,A,t),a(A,x),a(A,be),a(A,E),a(E,P),a(P,me),a(P,Y),a(Y,ee),a(P,fe),c(e,te,t),c(e,M,t),c(e,ae,t),c(e,W,t),c(e,le,t),c(e,U,t),c(e,ne,t),c(e,y,t),a(y,oe),a(y,ge),a(y,B),a(B,S),a(S,se),a(S,_e),a(S,ie),a(S,ke),a(S,m),a(m,ve),z(C,m,null),a(m,we),a(m,$e),a(m,Oe),a(m,re),a(m,Ae),a(m,ce),a(m,ye),a(m,Se),a(m,Te),a(m,de),a(m,Ce),a(B,qe),z(q,B,null),c(e,ue,t),c(e,F,t),c(e,he,t),c(e,T,t),a(T,L);for(let u=0;u<$.length;u+=1)$[u]&&$[u].m(L,null);a(T,Re),a(T,H);for(let u=0;u<w.length;u+=1)w[u]&&w[u].m(H,null);D=!0},p(e,[t]){(!D||t&1)&&i!==(i=e[0].name+"")&&pe(f,i);const u={};t&9&&(u.js=`
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
    `),v.$set(u),(!D||t&1)&&N!==(N=e[0].name+"")&&pe(ee,N),t&6&&(I=Q(e[2]),$=Ue($,t,Me,1,e,I,De,L,Ie,He,null,Le)),t&6&&(V=Q(e[2]),Qe(),w=Ue(w,t,We,1,e,V,Pe,H,Ze,Ve,null,Fe),ze())},i(e){if(!D){J(v.$$.fragment,e),J(C.$$.fragment,e),J(q.$$.fragment,e);for(let t=0;t<V.length;t+=1)J(w[t]);D=!0}},o(e){j(v.$$.fragment,e),j(C.$$.fragment,e),j(q.$$.fragment,e);for(let t=0;t<w.length;t+=1)j(w[t]);D=!1},d(e){e&&(r(l),r(d),r(b),r(k),r(O),r(R),r(G),r(A),r(te),r(M),r(ae),r(W),r(le),r(U),r(ne),r(y),r(ue),r(F),r(he),r(T)),Z(v,e),Z(C),Z(q);for(let t=0;t<$.length;t+=1)$[t].d();for(let t=0;t<w.length;t+=1)w[t].d()}}}function et(s,l,n){let i,{collection:f}=l,g=200,d=[];const b=k=>n(1,g=k.code);return s.$$set=k=>{"collection"in k&&n(0,f=k.collection)},s.$$.update=()=>{s.$$.dirty&1&&n(2,d=[{code:200,body:JSON.stringify({token:"JWT_AUTH_TOKEN",record:Be.dummyCollectionRecord(f),meta:{id:"abc123",name:"John Doe",username:"john.doe",email:"test@example.com",avatarURL:"https://example.com/avatar.png",accessToken:"...",refreshToken:"...",expiry:"2022-01-01 10:00:00.123Z",isNew:!1,rawUser:{}}},null,2)},{code:400,body:`
                {
                  "status": 400,
                  "message": "An error occurred while submitting the form.",
                  "data": {
                    "provider": {
                      "code": "validation_required",
                      "message": "Missing required value."
                    }
                  }
                }
            `}])},n(3,i=Be.getApiExampleUrl(Ke.baseURL)),[f,g,d,i,b]}class lt extends Je{constructor(l){super(),xe(this,l,et,Ye,Ee,{collection:0})}}export{lt as default};
