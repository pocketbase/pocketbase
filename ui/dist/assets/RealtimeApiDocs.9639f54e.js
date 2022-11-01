import{S as ee,i as oe,s as te,O as se,C as I,e as u,w as E,b as a,c as K,f as p,g as s,h as P,m as Q,x as ne,t as X,a as Z,o as n,d as x,L as ie,p as ce}from"./index.be8ffbe5.js";import{S as re}from"./SdkTabs.8f55857f.js";function le(t){var B,U,W,L,A,H,T,q;let i,m,c=t[0].name+"",b,d,k,f,S,v,w,r,_,$,O,g,y,h,D,l,R;return r=new re({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${t[1]}');

        ...

        // (Optionally) authenticate
        await pb.collection('users').authWithPassword('test@example.com', '123456');

        // Subscribe to changes in any record from the collection
        pb.collection('${(B=t[0])==null?void 0:B.name}').subscribe(function (e) {
            console.log(e.record);
        });

        // Subscribe to changes in a single record
        pb.collection('${(U=t[0])==null?void 0:U.name}').subscribeOne('RECORD_ID', function (e) {
            console.log(e.record);
        });

        // Unsubscribe
        pb.collection('${(W=t[0])==null?void 0:W.name}').unsubscribe() // remove all collection subscriptions
        pb.collection('${(L=t[0])==null?void 0:L.name}').unsubscribe('RECORD_ID') // remove only the record subscription
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${t[1]}');

        ...

        // (Optionally) authenticate
        await pb.collection('users').authWithPassword('test@example.com', '123456');

        // Subscribe to changes in any record from the collection
        pb.collection('${(A=t[0])==null?void 0:A.name}').subscribe((e) {
          print(e.record);
        });

        // Subscribe to changes in a single record
        pb.collection('${(H=t[0])==null?void 0:H.name}').subscribeOne('RECORD_ID', (e) {
          print(e.record);
        });

        // Unsubscribe
        pb.collection('${(T=t[0])==null?void 0:T.name}').unsubscribe() // remove all collection subscriptions
        pb.collection('${(q=t[0])==null?void 0:q.name}').unsubscribe('RECORD_ID') // remove only the record subscription
    `}}),l=new se({props:{content:JSON.stringify({action:"create",record:I.dummyCollectionRecord(t[0])},null,2).replace('"action": "create"','"action": "create" // create, update or delete')}}),{c(){i=u("h3"),m=E("Realtime ("),b=E(c),d=E(")"),k=a(),f=u("div"),f.innerHTML=`<p>Subscribe to realtime changes via Server-Sent Events (SSE).</p> 
    <p>Events are sent for <strong>create</strong>, <strong>update</strong>
        and <strong>delete</strong> record operations (see &quot;Event data format&quot; section below).</p>`,S=a(),v=u("div"),v.innerHTML=`<div class="icon"><i class="ri-information-line"></i></div> 
    <div class="contet"><p><strong>You could subscribe to a single record or to an entire collection.</strong></p> 
        <p>When you subscribe to a <strong>single record</strong>, the collection&#39;s
            <strong>ViewRule</strong> will be used to determine whether the subscriber has access to receive the
            event message.</p> 
        <p>When you subscribe to an <strong>entire collection</strong>, the collection&#39;s
            <strong>ListRule</strong> will be used to determine whether the subscriber has access to receive the
            event message.</p></div>`,w=a(),K(r.$$.fragment),_=a(),$=u("h6"),$.textContent="API details",O=a(),g=u("div"),g.innerHTML=`<strong class="label label-primary">SSE</strong> 
    <div class="content"><p>/api/realtime</p></div>`,y=a(),h=u("div"),h.textContent="Event data format",D=a(),K(l.$$.fragment),p(i,"class","m-b-sm"),p(f,"class","content txt-lg m-b-sm"),p(v,"class","alert alert-info m-t-10 m-b-sm"),p($,"class","m-b-xs"),p(g,"class","alert"),p(h,"class","section-title")},m(e,o){s(e,i,o),P(i,m),P(i,b),P(i,d),s(e,k,o),s(e,f,o),s(e,S,o),s(e,v,o),s(e,w,o),Q(r,e,o),s(e,_,o),s(e,$,o),s(e,O,o),s(e,g,o),s(e,y,o),s(e,h,o),s(e,D,o),Q(l,e,o),R=!0},p(e,[o]){var j,J,N,V,Y,z,F,G;(!R||o&1)&&c!==(c=e[0].name+"")&&ne(b,c);const C={};o&3&&(C.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[1]}');

        ...

        // (Optionally) authenticate
        await pb.collection('users').authWithPassword('test@example.com', '123456');

        // Subscribe to changes in any record from the collection
        pb.collection('${(j=e[0])==null?void 0:j.name}').subscribe(function (e) {
            console.log(e.record);
        });

        // Subscribe to changes in a single record
        pb.collection('${(J=e[0])==null?void 0:J.name}').subscribeOne('RECORD_ID', function (e) {
            console.log(e.record);
        });

        // Unsubscribe
        pb.collection('${(N=e[0])==null?void 0:N.name}').unsubscribe() // remove all collection subscriptions
        pb.collection('${(V=e[0])==null?void 0:V.name}').unsubscribe('RECORD_ID') // remove only the record subscription
    `),o&3&&(C.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[1]}');

        ...

        // (Optionally) authenticate
        await pb.collection('users').authWithPassword('test@example.com', '123456');

        // Subscribe to changes in any record from the collection
        pb.collection('${(Y=e[0])==null?void 0:Y.name}').subscribe((e) {
          print(e.record);
        });

        // Subscribe to changes in a single record
        pb.collection('${(z=e[0])==null?void 0:z.name}').subscribeOne('RECORD_ID', (e) {
          print(e.record);
        });

        // Unsubscribe
        pb.collection('${(F=e[0])==null?void 0:F.name}').unsubscribe() // remove all collection subscriptions
        pb.collection('${(G=e[0])==null?void 0:G.name}').unsubscribe('RECORD_ID') // remove only the record subscription
    `),r.$set(C);const M={};o&1&&(M.content=JSON.stringify({action:"create",record:I.dummyCollectionRecord(e[0])},null,2).replace('"action": "create"','"action": "create" // create, update or delete')),l.$set(M)},i(e){R||(X(r.$$.fragment,e),X(l.$$.fragment,e),R=!0)},o(e){Z(r.$$.fragment,e),Z(l.$$.fragment,e),R=!1},d(e){e&&n(i),e&&n(k),e&&n(f),e&&n(S),e&&n(v),e&&n(w),x(r,e),e&&n(_),e&&n($),e&&n(O),e&&n(g),e&&n(y),e&&n(h),e&&n(D),x(l,e)}}}function ae(t,i,m){let c,{collection:b=new ie}=i;return t.$$set=d=>{"collection"in d&&m(0,b=d.collection)},m(1,c=I.getApiExampleUrl(ce.baseUrl)),[b,c]}class pe extends ee{constructor(i){super(),oe(this,i,ae,le,te,{collection:0})}}export{pe as default};
