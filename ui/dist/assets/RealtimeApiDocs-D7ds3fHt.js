import{S as re,i as ae,s as be,V as pe,W as ue,J as P,h as s,d as se,t as ne,a as ie,I as me,l as n,n as y,m as ce,u as p,A as I,v as a,c as le,w as u,p as de}from"./index-TiFsHbkW.js";function he(o){var B,U,W,A,L,H,T,q,J,M,j,N;let i,m,c=o[0].name+"",b,d,k,h,D,f,_,l,S,$,w,g,C,v,E,r,R;return l=new pe({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${o[1]}');

        ...

        // (Optionally) authenticate
        await pb.collection('users').authWithPassword('test@example.com', '123456');

        // Subscribe to changes in any ${(B=o[0])==null?void 0:B.name} record
        pb.collection('${(U=o[0])==null?void 0:U.name}').subscribe('*', function (e) {
            console.log(e.action);
            console.log(e.record);
        }, { /* other options like: filter, expand, custom headers, etc. */ });

        // Subscribe to changes only in the specified record
        pb.collection('${(W=o[0])==null?void 0:W.name}').subscribe('RECORD_ID', function (e) {
            console.log(e.action);
            console.log(e.record);
        }, { /* other options like: filter, expand, custom headers, etc. */ });

        // Unsubscribe
        pb.collection('${(A=o[0])==null?void 0:A.name}').unsubscribe('RECORD_ID'); // remove all 'RECORD_ID' subscriptions
        pb.collection('${(L=o[0])==null?void 0:L.name}').unsubscribe('*'); // remove all '*' topic subscriptions
        pb.collection('${(H=o[0])==null?void 0:H.name}').unsubscribe(); // remove all subscriptions in the collection
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${o[1]}');

        ...

        // (Optionally) authenticate
        await pb.collection('users').authWithPassword('test@example.com', '123456');

        // Subscribe to changes in any ${(T=o[0])==null?void 0:T.name} record
        pb.collection('${(q=o[0])==null?void 0:q.name}').subscribe('*', (e) {
            print(e.action);
            print(e.record);
        }, /* other options like: filter, expand, custom headers, etc. */);

        // Subscribe to changes only in the specified record
        pb.collection('${(J=o[0])==null?void 0:J.name}').subscribe('RECORD_ID', (e) {
            print(e.action);
            print(e.record);
        }, /* other options like: filter, expand, custom headers, etc. */);

        // Unsubscribe
        pb.collection('${(M=o[0])==null?void 0:M.name}').unsubscribe('RECORD_ID'); // remove all 'RECORD_ID' subscriptions
        pb.collection('${(j=o[0])==null?void 0:j.name}').unsubscribe('*'); // remove all '*' topic subscriptions
        pb.collection('${(N=o[0])==null?void 0:N.name}').unsubscribe(); // remove all subscriptions in the collection
    `}}),r=new ue({props:{content:JSON.stringify({action:"create",record:P.dummyCollectionRecord(o[0])},null,2).replace('"action": "create"','"action": "create" // create, update or delete')}}),{c(){i=p("h3"),m=I("Realtime ("),b=I(c),d=I(")"),k=a(),h=p("div"),h.innerHTML=`<p>Subscribe to realtime changes via Server-Sent Events (SSE).</p> <p>Events are sent for <strong>create</strong>, <strong>update</strong>
        and <strong>delete</strong> record operations (see &quot;Event data format&quot; section below).</p>`,D=a(),f=p("div"),f.innerHTML=`<div class="icon"><i class="ri-information-line"></i></div> <div class="contet"><p><strong>You could subscribe to a single record or to an entire collection.</strong></p> <p>When you subscribe to a <strong>single record</strong>, the collection&#39;s
            <strong>ViewRule</strong> will be used to determine whether the subscriber has access to receive the
            event message.</p> <p>When you subscribe to an <strong>entire collection</strong>, the collection&#39;s
            <strong>ListRule</strong> will be used to determine whether the subscriber has access to receive the
            event message.</p></div>`,_=a(),le(l.$$.fragment),S=a(),$=p("h6"),$.textContent="API details",w=a(),g=p("div"),g.innerHTML='<strong class="label label-primary">SSE</strong> <div class="content"><p>/api/realtime</p></div>',C=a(),v=p("div"),v.textContent="Event data format",E=a(),le(r.$$.fragment),u(i,"class","m-b-sm"),u(h,"class","content txt-lg m-b-sm"),u(f,"class","alert alert-info m-t-10 m-b-sm"),u($,"class","m-b-xs"),u(g,"class","alert"),u(v,"class","section-title")},m(e,t){n(e,i,t),y(i,m),y(i,b),y(i,d),n(e,k,t),n(e,h,t),n(e,D,t),n(e,f,t),n(e,_,t),ce(l,e,t),n(e,S,t),n(e,$,t),n(e,w,t),n(e,g,t),n(e,C,t),n(e,v,t),n(e,E,t),ce(r,e,t),R=!0},p(e,[t]){var Y,z,F,G,K,Q,X,Z,x,ee,te,oe;(!R||t&1)&&c!==(c=e[0].name+"")&&me(b,c);const O={};t&3&&(O.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[1]}');

        ...

        // (Optionally) authenticate
        await pb.collection('users').authWithPassword('test@example.com', '123456');

        // Subscribe to changes in any ${(Y=e[0])==null?void 0:Y.name} record
        pb.collection('${(z=e[0])==null?void 0:z.name}').subscribe('*', function (e) {
            console.log(e.action);
            console.log(e.record);
        }, { /* other options like: filter, expand, custom headers, etc. */ });

        // Subscribe to changes only in the specified record
        pb.collection('${(F=e[0])==null?void 0:F.name}').subscribe('RECORD_ID', function (e) {
            console.log(e.action);
            console.log(e.record);
        }, { /* other options like: filter, expand, custom headers, etc. */ });

        // Unsubscribe
        pb.collection('${(G=e[0])==null?void 0:G.name}').unsubscribe('RECORD_ID'); // remove all 'RECORD_ID' subscriptions
        pb.collection('${(K=e[0])==null?void 0:K.name}').unsubscribe('*'); // remove all '*' topic subscriptions
        pb.collection('${(Q=e[0])==null?void 0:Q.name}').unsubscribe(); // remove all subscriptions in the collection
    `),t&3&&(O.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[1]}');

        ...

        // (Optionally) authenticate
        await pb.collection('users').authWithPassword('test@example.com', '123456');

        // Subscribe to changes in any ${(X=e[0])==null?void 0:X.name} record
        pb.collection('${(Z=e[0])==null?void 0:Z.name}').subscribe('*', (e) {
            print(e.action);
            print(e.record);
        }, /* other options like: filter, expand, custom headers, etc. */);

        // Subscribe to changes only in the specified record
        pb.collection('${(x=e[0])==null?void 0:x.name}').subscribe('RECORD_ID', (e) {
            print(e.action);
            print(e.record);
        }, /* other options like: filter, expand, custom headers, etc. */);

        // Unsubscribe
        pb.collection('${(ee=e[0])==null?void 0:ee.name}').unsubscribe('RECORD_ID'); // remove all 'RECORD_ID' subscriptions
        pb.collection('${(te=e[0])==null?void 0:te.name}').unsubscribe('*'); // remove all '*' topic subscriptions
        pb.collection('${(oe=e[0])==null?void 0:oe.name}').unsubscribe(); // remove all subscriptions in the collection
    `),l.$set(O);const V={};t&1&&(V.content=JSON.stringify({action:"create",record:P.dummyCollectionRecord(e[0])},null,2).replace('"action": "create"','"action": "create" // create, update or delete')),r.$set(V)},i(e){R||(ie(l.$$.fragment,e),ie(r.$$.fragment,e),R=!0)},o(e){ne(l.$$.fragment,e),ne(r.$$.fragment,e),R=!1},d(e){e&&(s(i),s(k),s(h),s(D),s(f),s(_),s(S),s($),s(w),s(g),s(C),s(v),s(E)),se(l,e),se(r,e)}}}function fe(o,i,m){let c,{collection:b}=i;return o.$$set=d=>{"collection"in d&&m(0,b=d.collection)},m(1,c=P.getApiExampleUrl(de.baseURL)),[b,c]}class ge extends re{constructor(i){super(),ae(this,i,fe,he,be,{collection:0})}}export{ge as default};
