import{S as re,i as ae,s as be,M as pe,C as P,e as p,w as y,b as a,c as oe,f as u,g as o,h as I,m as ne,x as ue,t as ie,a as ce,o as n,d as le,U as me,p as de}from"./index-a084d9d7.js";import{S as fe}from"./SdkTabs-ba0ec979.js";function $e(t){var U,B,W,A,H,L,M,T,q,j,J,N;let i,m,c=t[0].name+"",b,d,D,f,_,$,k,l,S,v,w,h,C,g,E,r,R;return l=new fe({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${t[1]}');

        ...

        // (Optionally) authenticate
        await pb.collection('users').authWithPassword('test@example.com', '123456');

        // Subscribe to changes in any ${(U=t[0])==null?void 0:U.name} record
        pb.collection('${(B=t[0])==null?void 0:B.name}').subscribe('*', function (e) {
            console.log(e.record);
        });

        // Subscribe to changes only in the specified record
        pb.collection('${(W=t[0])==null?void 0:W.name}').subscribe('RECORD_ID', function (e) {
            console.log(e.record);
        });

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

        // Subscribe to changes in any ${(M=t[0])==null?void 0:M.name} record
        pb.collection('${(T=t[0])==null?void 0:T.name}').subscribe('*', (e) {
            print(e.record);
        });

        // Subscribe to changes only in the specified record
        pb.collection('${(q=t[0])==null?void 0:q.name}').subscribe('RECORD_ID', (e) {
            print(e.record);
        });

        // Unsubscribe
        pb.collection('${(j=t[0])==null?void 0:j.name}').unsubscribe('RECORD_ID'); // remove all 'RECORD_ID' subscriptions
        pb.collection('${(J=t[0])==null?void 0:J.name}').unsubscribe('*'); // remove all '*' topic subscriptions
        pb.collection('${(N=t[0])==null?void 0:N.name}').unsubscribe(); // remove all subscriptions in the collection
    `}}),r=new pe({props:{content:JSON.stringify({action:"create",record:P.dummyCollectionRecord(t[0])},null,2).replace('"action": "create"','"action": "create" // create, update or delete')}}),{c(){i=p("h3"),m=y("Realtime ("),b=y(c),d=y(")"),D=a(),f=p("div"),f.innerHTML=`<p>Subscribe to realtime changes via Server-Sent Events (SSE).</p> <p>Events are sent for <strong>create</strong>, <strong>update</strong>
        and <strong>delete</strong> record operations (see &quot;Event data format&quot; section below).</p>`,_=a(),$=p("div"),$.innerHTML=`<div class="icon"><i class="ri-information-line"></i></div> <div class="contet"><p><strong>You could subscribe to a single record or to an entire collection.</strong></p> <p>When you subscribe to a <strong>single record</strong>, the collection&#39;s
            <strong>ViewRule</strong> will be used to determine whether the subscriber has access to receive the
            event message.</p> <p>When you subscribe to an <strong>entire collection</strong>, the collection&#39;s
            <strong>ListRule</strong> will be used to determine whether the subscriber has access to receive the
            event message.</p></div>`,k=a(),oe(l.$$.fragment),S=a(),v=p("h6"),v.textContent="API details",w=a(),h=p("div"),h.innerHTML='<strong class="label label-primary">SSE</strong> <div class="content"><p>/api/realtime</p></div>',C=a(),g=p("div"),g.textContent="Event data format",E=a(),oe(r.$$.fragment),u(i,"class","m-b-sm"),u(f,"class","content txt-lg m-b-sm"),u($,"class","alert alert-info m-t-10 m-b-sm"),u(v,"class","m-b-xs"),u(h,"class","alert"),u(g,"class","section-title")},m(e,s){o(e,i,s),I(i,m),I(i,b),I(i,d),o(e,D,s),o(e,f,s),o(e,_,s),o(e,$,s),o(e,k,s),ne(l,e,s),o(e,S,s),o(e,v,s),o(e,w,s),o(e,h,s),o(e,C,s),o(e,g,s),o(e,E,s),ne(r,e,s),R=!0},p(e,[s]){var Y,z,F,G,K,Q,X,Z,x,ee,se,te;(!R||s&1)&&c!==(c=e[0].name+"")&&ue(b,c);const O={};s&3&&(O.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[1]}');

        ...

        // (Optionally) authenticate
        await pb.collection('users').authWithPassword('test@example.com', '123456');

        // Subscribe to changes in any ${(Y=e[0])==null?void 0:Y.name} record
        pb.collection('${(z=e[0])==null?void 0:z.name}').subscribe('*', function (e) {
            console.log(e.record);
        });

        // Subscribe to changes only in the specified record
        pb.collection('${(F=e[0])==null?void 0:F.name}').subscribe('RECORD_ID', function (e) {
            console.log(e.record);
        });

        // Unsubscribe
        pb.collection('${(G=e[0])==null?void 0:G.name}').unsubscribe('RECORD_ID'); // remove all 'RECORD_ID' subscriptions
        pb.collection('${(K=e[0])==null?void 0:K.name}').unsubscribe('*'); // remove all '*' topic subscriptions
        pb.collection('${(Q=e[0])==null?void 0:Q.name}').unsubscribe(); // remove all subscriptions in the collection
    `),s&3&&(O.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[1]}');

        ...

        // (Optionally) authenticate
        await pb.collection('users').authWithPassword('test@example.com', '123456');

        // Subscribe to changes in any ${(X=e[0])==null?void 0:X.name} record
        pb.collection('${(Z=e[0])==null?void 0:Z.name}').subscribe('*', (e) {
            print(e.record);
        });

        // Subscribe to changes only in the specified record
        pb.collection('${(x=e[0])==null?void 0:x.name}').subscribe('RECORD_ID', (e) {
            print(e.record);
        });

        // Unsubscribe
        pb.collection('${(ee=e[0])==null?void 0:ee.name}').unsubscribe('RECORD_ID'); // remove all 'RECORD_ID' subscriptions
        pb.collection('${(se=e[0])==null?void 0:se.name}').unsubscribe('*'); // remove all '*' topic subscriptions
        pb.collection('${(te=e[0])==null?void 0:te.name}').unsubscribe(); // remove all subscriptions in the collection
    `),l.$set(O);const V={};s&1&&(V.content=JSON.stringify({action:"create",record:P.dummyCollectionRecord(e[0])},null,2).replace('"action": "create"','"action": "create" // create, update or delete')),r.$set(V)},i(e){R||(ie(l.$$.fragment,e),ie(r.$$.fragment,e),R=!0)},o(e){ce(l.$$.fragment,e),ce(r.$$.fragment,e),R=!1},d(e){e&&(n(i),n(D),n(f),n(_),n($),n(k),n(S),n(v),n(w),n(h),n(C),n(g),n(E)),le(l,e),le(r,e)}}}function ve(t,i,m){let c,{collection:b=new me}=i;return t.$$set=d=>{"collection"in d&&m(0,b=d.collection)},m(1,c=P.getApiExampleUrl(de.baseUrl)),[b,c]}class Re extends re{constructor(i){super(),ae(this,i,ve,$e,be,{collection:0})}}export{Re as default};
