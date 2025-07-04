import{S as St,i as At,s as Lt,a9 as Mt,aa as Ht,R as Y,h as d,d as Re,t as Z,a as x,I as jt,a0 as Pt,a1 as Nt,C as Ut,ab as Jt,D as zt,l as u,n as t,m as Te,E as Gt,G as Kt,u as a,A as _,v as i,c as Pe,w as b,J as Ft,p as Wt,k as ee,o as Qt}from"./index-CGZL2Fjn.js";function Bt(o,s,n){const c=o.slice();return c[6]=s[n],c}function Et(o,s,n){const c=o.slice();return c[6]=s[n],c}function Ot(o,s){let n,c,y;function f(){return s[5](s[6])}return{key:o,first:null,c(){n=a("button"),n.textContent=`${s[6].code} `,b(n,"class","tab-item"),ee(n,"active",s[1]===s[6].code),this.first=n},m(r,h){u(r,n,h),c||(y=Qt(n,"click",f),c=!0)},p(r,h){s=r,h&10&&ee(n,"active",s[1]===s[6].code)},d(r){r&&d(n),c=!1,y()}}}function It(o,s){let n,c,y,f;return c=new Ht({props:{content:s[6].body}}),{key:o,first:null,c(){n=a("div"),Pe(c.$$.fragment),y=i(),b(n,"class","tab-item"),ee(n,"active",s[1]===s[6].code),this.first=n},m(r,h){u(r,n,h),Te(c,n,null),t(n,y),f=!0},p(r,h){s=r,(!f||h&10)&&ee(n,"active",s[1]===s[6].code)},i(r){f||(x(c.$$.fragment,r),f=!0)},o(r){Z(c.$$.fragment,r),f=!1},d(r){r&&d(n),Re(c)}}}function Vt(o){var pt,mt,bt,ht,ft,_t,yt,kt;let s,n,c=o[0].name+"",y,f,r,h,F,g,U,Fe,P,B,Be,E,Ee,Oe,te,le,q,ae,O,oe,I,se,H,ne,J,ie,w,ce,Ie,re,S,z,He,k,G,Se,de,Ae,C,K,Le,ue,Me,W,je,pe,Ne,D,Ue,me,Je,ze,Ge,Q,Ke,V,We,be,Qe,he,Ve,fe,Xe,p,_e,Ye,ye,Ze,ke,xe,$e,et,ge,tt,ve,lt,at,ot,Ce,st,R,De,A,qe,T,L,v=[],nt=new Map,it,M,$=[],ct=new Map,j,we,rt;q=new Mt({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${o[2]}');

        ...

        const batch = pb.createBatch();

        batch.collection('${(pt=o[0])==null?void 0:pt.name}').create({ ... });
        batch.collection('${(mt=o[0])==null?void 0:mt.name}').update('RECORD_ID', { ... });
        batch.collection('${(bt=o[0])==null?void 0:bt.name}').delete('RECORD_ID');
        batch.collection('${(ht=o[0])==null?void 0:ht.name}').upsert({ ... });

        const result = await batch.send();
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${o[2]}');

        ...

        final batch = pb.createBatch();

        batch.collection('${(ft=o[0])==null?void 0:ft.name}').create(body: { ... });
        batch.collection('${(_t=o[0])==null?void 0:_t.name}').update('RECORD_ID', body: { ... });
        batch.collection('${(yt=o[0])==null?void 0:yt.name}').delete('RECORD_ID');
        batch.collection('${(kt=o[0])==null?void 0:kt.name}').upsert(body: { ... });

        final result = await batch.send();
    `}}),R=new Ht({props:{language:"javascript",content:`
                            const formData = new FormData();

                            formData.append("@jsonPayload", JSON.stringify({
                                requests: [
                                    {
                                        method: "POST",
                                        url: "/api/collections/${o[0].name}/records?fields=id",
                                        body: { someField: "test1" }
                                    },
                                    {
                                        method: "PATCH",
                                        url: "/api/collections/${o[0].name}/records/RECORD_ID",
                                        body: { someField: "test2" }
                                    }
                                ]
                            }))

                            // file for the first request
                            formData.append("requests.0.someFileField", new File(...))

                            // file for the second request
                            formData.append("requests.1.someFileField", new File(...))
                        `}});let X=Y(o[3]);const dt=e=>e[6].code;for(let e=0;e<X.length;e+=1){let l=Et(o,X,e),m=dt(l);nt.set(m,v[e]=Ot(m,l))}let N=Y(o[3]);const ut=e=>e[6].code;for(let e=0;e<N.length;e+=1){let l=Bt(o,N,e),m=ut(l);ct.set(m,$[e]=It(m,l))}return{c(){s=a("h3"),n=_("Batch create/update/upsert/delete ("),y=_(c),f=_(")"),r=i(),h=a("div"),h.innerHTML="<p>Batch and transactional create/update/upsert/delete of multiple records in a single request.</p>",F=i(),g=a("div"),U=a("div"),U.innerHTML='<i class="ri-error-warning-line"></i>',Fe=i(),P=a("div"),B=a("p"),Be=_(`The batch Web API need to be explicitly enabled and configured from the
            `),E=a("a"),E.textContent="Dashboard settings",Ee=_("."),Oe=i(),te=a("p"),te.textContent=`Because this endpoint process the requests in a single transaction it could degrade the
            performance of your application if not used with proper care and configuration (e.g. too large
            allowed execution timeout, large body size limit, etc.).`,le=i(),Pe(q.$$.fragment),ae=i(),O=a("h6"),O.textContent="API details",oe=i(),I=a("div"),I.innerHTML='<strong class="label label-primary">POST</strong> <div class="content">/api/batch</div>',se=i(),H=a("div"),H.textContent="Body Parameters",ne=i(),J=a("p"),J.innerHTML=`Body parameters could be sent as <em>application/json</em> or <em>multipart/form-data</em>.
    <br/>
    File upload is supported only via <em>multipart/form-data</em> (see below for more details).`,ie=i(),w=a("table"),ce=a("thead"),ce.innerHTML='<tr><th>Param</th> <th width="80%">Description</th></tr>',Ie=i(),re=a("tbody"),S=a("tr"),z=a("td"),z.innerHTML='<div class="flex txt-nowrap"><span class="label label-success">Required</span> <span>requests</span></div>',He=i(),k=a("td"),G=a("span"),G.textContent="Array<Request>",Se=_(` - List of the requests to process.

                `),de=a("p"),de.textContent="The supported batch request actions are:",Ae=i(),C=a("ul"),K=a("li"),Le=_("record create - "),ue=a("code"),ue.textContent="POST /api/collections/{collection}/records",Me=i(),W=a("li"),je=_(`record update -
                        `),pe=a("code"),pe.textContent="PATCH /api/collections/{collection}/records/{id}",Ne=i(),D=a("li"),Ue=_("record upsert - "),me=a("code"),me.textContent="PUT /api/collections/{collection}/records",Je=i(),ze=a("br"),Ge=i(),Q=a("small"),Q.innerHTML='(the body must have <code class="txt-sm">id</code> field)',Ke=i(),V=a("li"),We=_(`record delete -
                        `),be=a("code"),be.textContent="DELETE /api/collections/{collection}/records/{id}",Qe=i(),he=a("p"),he.textContent="Each batch Request element have the following properties:",Ve=i(),fe=a("ul"),fe.innerHTML=`<li><code>url path</code> <em>(could include query parameters)</em></li> <li><code>method</code> <em>(GET, POST, PUT, PATCH, DELETE)</em></li> <li><code>headers</code> <br/> <em>(custom per-request <code>Authorization</code> header is not supported at the moment,
                            aka. all batch requests have the same auth state)</em></li> <li><code>body</code></li>`,Xe=i(),p=a("p"),_e=a("strong"),_e.textContent="NB!",Ye=_(` When the batch request is send as
                    `),ye=a("code"),ye.textContent="multipart/form-data",Ze=_(`, the regular batch action fields are expected to be
                    submitted as serailized json under the `),ke=a("code"),ke.textContent="@jsonPayload",xe=_(` field and file keys need
                    to follow the pattern `),$e=a("code"),$e.textContent="requests.N.fileField",et=_(` or
                    `),ge=a("code"),ge.textContent="requests[N].fileField",tt=i(),ve=a("em"),ve.textContent=`(this is usually handled transparently by the SDKs when their specific object notation
                        is used)
                    `,lt=_(`.
                    `),at=a("br"),ot=_(`
                    If you don't use the SDKs or prefer manually to construct the `),Ce=a("code"),Ce.textContent="FormData",st=_(`
                    body, then it could look something like:
                    `),Pe(R.$$.fragment),De=i(),A=a("div"),A.textContent="Responses",qe=i(),T=a("div"),L=a("div");for(let e=0;e<v.length;e+=1)v[e].c();it=i(),M=a("div");for(let e=0;e<$.length;e+=1)$[e].c();b(s,"class","m-b-sm"),b(h,"class","content txt-lg m-b-sm"),b(U,"class","icon"),b(E,"href","/settings"),b(P,"class","content"),b(g,"class","alert alert-warning"),b(O,"class","m-b-xs"),b(I,"class","api-route alert alert-success"),b(H,"class","section-title"),b(z,"valign","top"),b(G,"class","label"),b(Q,"class","txt-hint"),b(w,"class","table-compact table-border m-t-xs m-b-base"),b(A,"class","section-title"),b(L,"class","tabs-header compact combined left"),b(M,"class","tabs-content"),b(T,"class","tabs")},m(e,l){u(e,s,l),t(s,n),t(s,y),t(s,f),u(e,r,l),u(e,h,l),u(e,F,l),u(e,g,l),t(g,U),t(g,Fe),t(g,P),t(P,B),t(B,Be),t(B,E),t(B,Ee),t(P,Oe),t(P,te),u(e,le,l),Te(q,e,l),u(e,ae,l),u(e,O,l),u(e,oe,l),u(e,I,l),u(e,se,l),u(e,H,l),u(e,ne,l),u(e,J,l),u(e,ie,l),u(e,w,l),t(w,ce),t(w,Ie),t(w,re),t(re,S),t(S,z),t(S,He),t(S,k),t(k,G),t(k,Se),t(k,de),t(k,Ae),t(k,C),t(C,K),t(K,Le),t(K,ue),t(C,Me),t(C,W),t(W,je),t(W,pe),t(C,Ne),t(C,D),t(D,Ue),t(D,me),t(D,Je),t(D,ze),t(D,Ge),t(D,Q),t(C,Ke),t(C,V),t(V,We),t(V,be),t(k,Qe),t(k,he),t(k,Ve),t(k,fe),t(k,Xe),t(k,p),t(p,_e),t(p,Ye),t(p,ye),t(p,Ze),t(p,ke),t(p,xe),t(p,$e),t(p,et),t(p,ge),t(p,tt),t(p,ve),t(p,lt),t(p,at),t(p,ot),t(p,Ce),t(p,st),Te(R,p,null),u(e,De,l),u(e,A,l),u(e,qe,l),u(e,T,l),t(T,L);for(let m=0;m<v.length;m+=1)v[m]&&v[m].m(L,null);t(T,it),t(T,M);for(let m=0;m<$.length;m+=1)$[m]&&$[m].m(M,null);j=!0,we||(rt=Gt(Kt.call(null,E)),we=!0)},p(e,[l]){var gt,vt,Ct,Dt,qt,wt,Rt,Tt;(!j||l&1)&&c!==(c=e[0].name+"")&&jt(y,c);const m={};l&5&&(m.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[2]}');

        ...

        const batch = pb.createBatch();

        batch.collection('${(gt=e[0])==null?void 0:gt.name}').create({ ... });
        batch.collection('${(vt=e[0])==null?void 0:vt.name}').update('RECORD_ID', { ... });
        batch.collection('${(Ct=e[0])==null?void 0:Ct.name}').delete('RECORD_ID');
        batch.collection('${(Dt=e[0])==null?void 0:Dt.name}').upsert({ ... });

        const result = await batch.send();
    `),l&5&&(m.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[2]}');

        ...

        final batch = pb.createBatch();

        batch.collection('${(qt=e[0])==null?void 0:qt.name}').create(body: { ... });
        batch.collection('${(wt=e[0])==null?void 0:wt.name}').update('RECORD_ID', body: { ... });
        batch.collection('${(Rt=e[0])==null?void 0:Rt.name}').delete('RECORD_ID');
        batch.collection('${(Tt=e[0])==null?void 0:Tt.name}').upsert(body: { ... });

        final result = await batch.send();
    `),q.$set(m);const $t={};l&1&&($t.content=`
                            const formData = new FormData();

                            formData.append("@jsonPayload", JSON.stringify({
                                requests: [
                                    {
                                        method: "POST",
                                        url: "/api/collections/${e[0].name}/records?fields=id",
                                        body: { someField: "test1" }
                                    },
                                    {
                                        method: "PATCH",
                                        url: "/api/collections/${e[0].name}/records/RECORD_ID",
                                        body: { someField: "test2" }
                                    }
                                ]
                            }))

                            // file for the first request
                            formData.append("requests.0.someFileField", new File(...))

                            // file for the second request
                            formData.append("requests.1.someFileField", new File(...))
                        `),R.$set($t),l&10&&(X=Y(e[3]),v=Pt(v,l,dt,1,e,X,nt,L,Nt,Ot,null,Et)),l&10&&(N=Y(e[3]),Ut(),$=Pt($,l,ut,1,e,N,ct,M,Jt,It,null,Bt),zt())},i(e){if(!j){x(q.$$.fragment,e),x(R.$$.fragment,e);for(let l=0;l<N.length;l+=1)x($[l]);j=!0}},o(e){Z(q.$$.fragment,e),Z(R.$$.fragment,e);for(let l=0;l<$.length;l+=1)Z($[l]);j=!1},d(e){e&&(d(s),d(r),d(h),d(F),d(g),d(le),d(ae),d(O),d(oe),d(I),d(se),d(H),d(ne),d(J),d(ie),d(w),d(De),d(A),d(qe),d(T)),Re(q,e),Re(R);for(let l=0;l<v.length;l+=1)v[l].d();for(let l=0;l<$.length;l+=1)$[l].d();we=!1,rt()}}}function Xt(o,s,n){let c,y,{collection:f}=s,r=200,h=[];const F=g=>n(1,r=g.code);return o.$$set=g=>{"collection"in g&&n(0,f=g.collection)},o.$$.update=()=>{o.$$.dirty&1&&n(4,y=Ft.dummyCollectionRecord(f)),o.$$.dirty&17&&f!=null&&f.id&&(h.push({code:200,body:JSON.stringify([{status:200,body:y},{status:200,body:Object.assign({},y,{id:y.id+"2"})}],null,2)}),h.push({code:400,body:`
                {
                  "status": 400,
                  "message": "Batch transaction failed.",
                  "data": {
                    "requests": {
                      "1": {
                        "code": "batch_request_failed",
                        "message": "Batch request failed.",
                        "response": {
                          "status": 400,
                          "message": "Failed to create record.",
                          "data": {
                            "id": {
                              "code": "validation_min_text_constraint",
                              "message": "Must be at least 3 character(s).",
                              "params": { "min": 3 }
                            }
                          }
                        }
                      }
                    }
                  }
                }
            `}),h.push({code:403,body:`
                {
                  "status": 403,
                  "message": "Batch requests are not allowed.",
                  "data": {}
                }
            `}))},n(2,c=Ft.getApiExampleUrl(Wt.baseURL)),[f,r,c,h,y,F]}class Zt extends St{constructor(s){super(),At(this,s,Xt,Vt,Lt,{collection:0})}}export{Zt as default};
