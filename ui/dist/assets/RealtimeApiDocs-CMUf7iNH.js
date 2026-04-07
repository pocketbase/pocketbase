import{S as re,i as ae,s as be,N as pe,C as P,b as s,d as se,t as ne,a as ie,r as ue,f as n,h as y,m as ce,n as p,u as I,k as a,c as le,o as u,p as me}from"./index--FBvE7un.js";import{S as de}from"./SdkTabs-E-5sYnXA.js";function he(t){var B,U,W,A,H,L,T,q,M,N,j,J;let i,m,c=t[0].name+"",b,d,R,h,D,f,_,l,S,$,C,g,E,v,w,r,k;return l=new de({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${t[1]}');

        ...

        // (Optionally) authenticate
        await pb.collection('users').authWithPassword('test@example.com', '123456');

        // Subscribe to changes in any ${(B=t[0])==null?void 0:B.name} record
        pb.collection('${(U=t[0])==null?void 0:U.name}').subscribe('*', function (e) {
            console.log(e.action);
            console.log(e.record);
        }, { /* other options like expand, custom headers, etc. */ });

        // Subscribe to changes only in the specified record
        pb.collection('${(W=t[0])==null?void 0:W.name}').subscribe('RECORD_ID', function (e) {
            console.log(e.action);
            console.log(e.record);
        }, { /* other options like expand, custom headers, etc. */ });

        // Unsubscribe
        pb.collection('${(A=t[0])==null?void 0:A.name}').unsubscribe('RECORD_ID'); // remove all 'RECORD_ID' subscriptions
        pb.collection('${(H=t[0])==null?void 0:H.name}').unsubscribe('*'); // remove all '*' topic subscriptions
        pb.collection('${(L=t[0])==null?void 0:L.name}').unsubscribe(); // remove all subscriptions in the collection
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${t[1]}');

        ...

        // (Optionally) authenticate
        await pb.collection('users').authWithPassword('test@example.com', '123456');

        // Subscribe to changes in any ${(T=t[0])==null?void 0:T.name} record
        pb.collection('${(q=t[0])==null?void 0:q.name}').subscribe('*', (e) {
            print(e.action);
            print(e.record);
        }, /* other options like expand, custom headers, etc. */);

        // Subscribe to changes only in the specified record
        pb.collection('${(M=t[0])==null?void 0:M.name}').subscribe('RECORD_ID', (e) {
            print(e.action);
            print(e.record);
        }, /* other options like expand, custom headers, etc. */);

        // Unsubscribe
        pb.collection('${(N=t[0])==null?void 0:N.name}').unsubscribe('RECORD_ID'); // remove all 'RECORD_ID' subscriptions
        pb.collection('${(j=t[0])==null?void 0:j.name}').unsubscribe('*'); // remove all '*' topic subscriptions
        pb.collection('${(J=t[0])==null?void 0:J.name}').unsubscribe(); // remove all subscriptions in the collection
    `}}),r=new pe({props:{content:JSON.stringify({action:"create",record:P.dummyCollectionRecord(t[0])},null,2).replace('"action": "create"','"action": "create" // create, update or delete')}}),{c(){i=p("h3"),m=I("Realtime ("),b=I(c),d=I(")"),R=a(),h=p("div"),h.innerHTML=`<p>Subscribe to realtime changes via Server-Sent Events (SSE).</p> <p>Events are sent for <strong>create</strong>, <strong>update</strong>
        and <strong>delete</strong> record operations (see &quot;Event data format&quot; section below).</p>`,D=a(),f=p("div"),f.innerHTML=`<div class="icon"><i class="ri-information-line"></i></div> <div class="contet"><p><strong>You could subscribe to a single record or to an entire collection.</strong></p> <p>When you subscribe to a <strong>single record</strong>, the collection&#39;s
            <strong>ViewRule</strong> will be used to determine whether the subscriber has access to receive the
            event message.</p> <p>When you subscribe to an <strong>entire collection</strong>, the collection&#39;s
            <strong>ListRule</strong> will be used to determine whether the subscriber has access to receive the
            event message.</p></div>`,_=a(),le(l.$$.fragment),S=a(),$=p("h6"),$.textContent="API details",C=a(),g=p("div"),g.innerHTML='<strong class="label label-primary">SSE</strong> <div class="content"><p>/api/realtime</p></div>',E=a(),v=p("div"),v.textContent="Event data format",w=a(),le(r.$$.fragment),u(i,"class","m-b-sm"),u(h,"class","content txt-lg m-b-sm"),u(f,"class","alert alert-info m-t-10 m-b-sm"),u($,"class","m-b-xs"),u(g,"class","alert"),u(v,"class","section-title")},m(e,o){n(e,i,o),y(i,m),y(i,b),y(i,d),n(e,R,o),n(e,h,o),n(e,D,o),n(e,f,o),n(e,_,o),ce(l,e,o),n(e,S,o),n(e,$,o),n(e,C,o),n(e,g,o),n(e,E,o),n(e,v,o),n(e,w,o),ce(r,e,o),k=!0},p(e,[o]){var Y,z,F,G,K,Q,X,Z,x,ee,oe,te;(!k||o&1)&&c!==(c=e[0].name+"")&&ue(b,c);const O={};o&3&&(O.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[1]}');

        ...

        // (Optionally) authenticate
        await pb.collection('users').authWithPassword('test@example.com', '123456');

        // Subscribe to changes in any ${(Y=e[0])==null?void 0:Y.name} record
        pb.collection('${(z=e[0])==null?void 0:z.name}').subscribe('*', function (e) {
            console.log(e.action);
            console.log(e.record);
        }, { /* other options like expand, custom headers, etc. */ });

        // Subscribe to changes only in the specified record
        pb.collection('${(F=e[0])==null?void 0:F.name}').subscribe('RECORD_ID', function (e) {
            console.log(e.action);
            console.log(e.record);
        }, { /* other options like expand, custom headers, etc. */ });

        // Unsubscribe
        pb.collection('${(G=e[0])==null?void 0:G.name}').unsubscribe('RECORD_ID'); // remove all 'RECORD_ID' subscriptions
        pb.collection('${(K=e[0])==null?void 0:K.name}').unsubscribe('*'); // remove all '*' topic subscriptions
        pb.collection('${(Q=e[0])==null?void 0:Q.name}').unsubscribe(); // remove all subscriptions in the collection
    `),o&3&&(O.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[1]}');

        ...

        // (Optionally) authenticate
        await pb.collection('users').authWithPassword('test@example.com', '123456');

        // Subscribe to changes in any ${(X=e[0])==null?void 0:X.name} record
        pb.collection('${(Z=e[0])==null?void 0:Z.name}').subscribe('*', (e) {
            print(e.action);
            print(e.record);
        }, /* other options like expand, custom headers, etc. */);

        // Subscribe to changes only in the specified record
        pb.collection('${(x=e[0])==null?void 0:x.name}').subscribe('RECORD_ID', (e) {
            print(e.action);
            print(e.record);
        }, /* other options like expand, custom headers, etc. */);

        // Unsubscribe
        pb.collection('${(ee=e[0])==null?void 0:ee.name}').unsubscribe('RECORD_ID'); // remove all 'RECORD_ID' subscriptions
        pb.collection('${(oe=e[0])==null?void 0:oe.name}').unsubscribe('*'); // remove all '*' topic subscriptions
        pb.collection('${(te=e[0])==null?void 0:te.name}').unsubscribe(); // remove all subscriptions in the collection
    `),l.$set(O);const V={};o&1&&(V.content=JSON.stringify({action:"create",record:P.dummyCollectionRecord(e[0])},null,2).replace('"action": "create"','"action": "create" // create, update or delete')),r.$set(V)},i(e){k||(ie(l.$$.fragment,e),ie(r.$$.fragment,e),k=!0)},o(e){ne(l.$$.fragment,e),ne(r.$$.fragment,e),k=!1},d(e){e&&(s(i),s(R),s(h),s(D),s(f),s(_),s(S),s($),s(C),s(g),s(E),s(v),s(w)),se(l,e),se(r,e)}}}function fe(t,i,m){let c,{collection:b}=i;return t.$$set=d=>{"collection"in d&&m(0,b=d.collection)},m(1,c=P.getApiExampleUrl(me.baseUrl)),[b,c]}class ve extends re{constructor(i){super(),ae(this,i,fe,he,be,{collection:0})}}export{ve as default};
