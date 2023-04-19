document.addEventListener("DOMContentLoaded", function () {
    $('#grad').click(function () { 
        {
            $.ajax({
                url: "resultado.php", 
                type: "POST",
                data: $("#graduaciones").serialize(), 
                success: function (resultado) {
                    $("#okgrad").html(resultado);  

                }
            });
        }
    })



    $('#tbl').DataTable();
    $(".confirmar").submit(function (e) {
        e.preventDefault();
        Swal.fire({
            title: 'Esta seguro de eliminar?',
            icon: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#3085d6',
            cancelButtonColor: '#d33',
            confirmButtonText: 'SI, Eliminar!'
        }).then((result) => {
            if (result.isConfirmed) {
                this.submit();
            }
        })
    })
    $("#nom_cliente").autocomplete({
        minLength: 3,
        source: function (request, response) {
            $.ajax({
                url: "ajax.php",
                dataType: "json",
                data: {
                    q: request.term
                },
                success: function (data) {
                    response(data);
                }
            });
        },
        select: function (event, ui) {
            $("#idcliente").val(ui.item.id);
            $("#nom_cliente").val(ui.item.label);
            $("#tel_cliente").val(ui.item.telefono);
            $("#dir_cliente").val(ui.item.direccion);
            $("#obrasocial").val(ui.item.obrasocial);
        }
    })
    $("#producto").autocomplete({
        minLength: 3,
        source: function (request, response) {
            $.ajax({
                url: "ajax.php",
                dataType: "json",
                data: {
                    pro: request.term
                },
                success: function (data) {
                    response(data);
                }
            });
        },
        select: function (event, ui) {
            $("#producto").val(ui.item.value);
            setTimeout(
                function () {
                    e = jQuery.Event("keypress");
                    e.which = 13;
                    registrarDetalle(e, ui.item.id, 1, ui.item.precio);
                }
            )
        } 
    })
    $('#btn_generar').click(function (e) {
        e.preventDefault();
        var rows = $('#tblDetalle tr').length;
        if (rows > 2) {
            var abona = $('#abona').val();
            var action = 'procesarVenta';
            var id = $('#idcliente').val();            
            var resto = $('#resto').val();
            var descuento = $('#porc').val();
            var metodo_pago = $('input[name=pago]:checked').val();

            var obrasocial = $('#obra_social').val();
            if (abona == "" || abona == null) {
                Swal.fire({
                    position: 'top-end',
                    icon: 'error',
                    title: 'El campo abona no puede estar vacio',
                    showConfirmButton: false,
                    timer: 2000
                });
                return;
            }
            $.ajax({
                url: 'ajax.php',
                async: true,
                data: {
                    procesarVenta: action,
                    id: id,
                    abona : abona,
                    resto : resto,
                    descuento : descuento,
                    obrasocial : obrasocial,
                    metodo_pago : metodo_pago 
                
                },
                success: function (response) {
                    const res = JSON.parse(response);
                    if (response != 'error') {
                        Swal.fire({
                            position: 'top-end',
                            icon: 'success',
                            title: 'Venta Generada',
                            showConfirmButton: false,
                            timer: 2000
                        })
                        setTimeout(() => {
                            generarPDF(res.id_cliente, res.id_venta);
                            location.reload();
                        }, 300);
                    } else {
                        Swal.fire({
                            position: 'top-end',
                            icon: 'error',
                            title: 'Error al generar la venta',
                            showConfirmButton: false,
                            timer: 2000
                        })
                    }
                
                },
               
                error: function (error) {

                }
            });
        } else {
            Swal.fire({
                position: 'top-end',
                icon: 'warning',
                title: 'No hay producto para generar la venta',
                showConfirmButton: false,
                timer: 2000
            })
        }
    });
    if (document.getElementById("detalle_venta")) {
        listar();
    }

    document.querySelector("#borrar_grad").addEventListener("click", function () {
        {
            $.ajax({
                url: "borrar_grad.php",
                type: "POST",
                data: $("#borrar_grad").serialize(),
                success: function (resultado) {
                    $("#okgrad").html(resultado);
    
                }
            });
        }
    })
    
    
    
    

})
document.querySelector("#guardar_cristal").addEventListener("click", function () {
    {
        $.ajax({
            url: "colocar_cristal.php",
            type: "POST",
            data: $("#form_cristal").serialize(),
            success: function (resultado) {
                $("#div_cristal").html(resultado);
            }
        });
    }
})

document.querySelector("#buscar_venta").addEventListener("click", function () {
    {
        $.ajax({
            url: "postpagos.php",
            type: "POST",
            data: $("#form_venta").serialize(),
            success: function (resultado) {
                $("#div_venta").html(resultado);

            }
        });
    }
})

document.querySelector("#anular_venta").addEventListener("click", function () {
    var idventa = $('#idanular').val();
    if(idventa == ""){
         Swal.fire({
            position: 'top-mid',
            icon: 'error',
            title: 'Complete Id Venta',
            showConfirmButton: false,
            timer: 2000
        });
        return;
    }

    Swal.fire({
        position: 'top-mid',
        icon: 'success',
        title: '',
        text: '¿Desea Eliminar la Venta?',
        showConfirmButton: true,
        showCancelButton: true,
        confirmButtonColor: '#3085d6',
        cancelButtonColor: '#d33',
        confirmButtonText: 'Si, Eliminar!'    
    }).then((result) => {
        if (result.isConfirmed) {
            // La respuesta del usuario es "Aceptar"
            {
                $.ajax({
                    url: "anular.php",
                    type: "POST",
                    data: $("#form_anular").serialize(),
                    success: function (resultado) {
                        $("#div_anular").html(resultado);
                    }
                });
            }
        } else if (result.isDenied || result.isDismissed) {
            // La respuesta del usuario es "Cancelar"
            Swal.fire('No se ha eliminado la venta', '', 'info')
        }
    });
});
    

function listar() {
    let html = '';
    let detalle = 'detalle';
    $.ajax({
        url: "ajax.php",
        dataType: "json",
        data: {
            detalle: detalle
        },
        success: function (response) {
            response.forEach(row => {
                html += `<tr>
                <td>${row['id']}</td>
                <td>${row['descripcion']}</td>
                <td>${row['cantidad']}</td>
                <td>${row['precio_venta']}</td>
                <td>${row['sub_total']}</td>
                <td><button class="btn btn-danger" type="button" onclick="deleteDetalle(${row['id']})">
                <i class="fas fa-trash-alt"></i></button></td>
                </tr>`;
            });
            document.querySelector("#detalle_venta").innerHTML = html;
            calcular();

        }
    });
}


function registrarDetalle(e, id, cant, precio) {
    if (document.getElementById('producto').value != '') {
        if (e.which == 13) {
            if (id != null) {
                let action = 'regDetalle';
                $.ajax({
                    url: "ajax.php",
                    type: 'POST',
                    dataType: "json",
                    data: {
                        id: id,
                        cant: cant,
                        action: action,
                        precio: precio
                    },
                    success: function (response) {
                        if (response == 'registrado') {
                            Swal.fire({
                                position: 'top-end',
                                icon: 'success',
                                title: 'Producto Ingresado',
                                showConfirmButton: false,
                                timer: 2000
                            })
                            document.querySelector("#producto").value = '';
                            document.querySelector("#producto").focus();
                            listar();
                        } else if (response == 'actualizado') {
                            Swal.fire({
                                position: 'top-end',
                                icon: 'success',
                                title: 'Producto Actualizado',
                                showConfirmButton: false,
                                timer: 2000
                            })
                            document.querySelector("#producto").value = '';
                            document.querySelector("#producto").focus();
                            listar();
                        } else {
                            Swal.fire({
                                position: 'top-end',
                                icon: 'error',
                                title: 'Error al ingresar el producto',
                                showConfirmButton: false,
                                timer: 2000
                            })
                        }
                    }
                });
            }
        }
    }
}
function deleteDetalle(id) {
    let detalle = 'Eliminar'
    $.ajax({
        url: "ajax.php",
        data: {
            id: id,
            delete_detalle: detalle
        },
        success: function (response) {
            console.log(response);
            if (response == 'restado') {
                Swal.fire({
                    position: 'top-end',
                    icon: 'success',
                    title: 'Producto Descontado',
                    showConfirmButton: false,
                    timer: 2000
                })
                document.querySelector("#producto").value = '';
                document.querySelector("#producto").focus();
                listar();
            } else if (response == 'ok') {
                Swal.fire({
                    position: 'top-end',
                    icon: 'success',
                    title: 'Producto Eliminado',
                    showConfirmButton: false,
                    timer: 2000
                })
                document.querySelector("#producto").value = '';
                document.querySelector("#producto").focus();
                listar();
            } else {
                Swal.fire({
                    position: 'top-end',
                    icon: 'error',
                    title: 'Error al eliminar el producto',
                    showConfirmButton: false,
                    timer: 2000
                })
            }
        }
    });
}


function calcular() {
    var total = 0;
    // obtenemos todas las filas del tbody
    var filas = document.querySelectorAll("#tblDetalle tbody tr");

    // recorremos cada una de las filas
    filas.forEach(function (e) {
        // obtenemos las columnas de cada fila
        var columnas = e.querySelectorAll("td");

        // obtenemos los valores de la cantidad y importe
        var importe = parseFloat(columnas[4].textContent);

        total += importe;
    });

    // mostramos la suma total
    var filas = document.querySelectorAll("#tblDetalle tfoot tr td");
    filas[1].textContent = total.toFixed(2);

    document.querySelector("#btn_parcial").addEventListener("click", function (total) {
        var abona = document.getElementById('abona');
        if (!abona.value || abona.value > total.value) {
            Swal.fire({
                position: 'top-end',
                icon: 'error',
                title: 'Error en campo abona, revise',
                showConfirmButton: false,
                timer: 2000
            });
            return;
        }else{
        var descuento = document.getElementById('porc');
        var obrasocial = document.getElementById('obra_social');
        if(!obrasocial.value || obrasocial.value <= 0 ){
            obrasocial.value = 0;
        }
        if(obrasocial.value > total.value){
            Swal.fire({
                position: 'top-end',
                icon: 'error',
                title: 'Obra Social es mayor que Total',
                showConfirmButton: false,
                timer: 2000
            });
            return; 
        }
        if((obrasocial.value + abona.value) > total.value){
            Swal.fire({
                position: 'top-end',
                icon: 'error',
                title: 'Obra Social más Abona es mayor que Total',
                showConfirmButton: false,
                timer: 2000
            });
            return; 
        }
        var resto = document.getElementById('resto');
        var dto = descuento.value;
        total = (total - obrasocial.value) * dto;
        var filas = document.querySelectorAll("#tblDetalle tfoot tr td");
        filas[1].textContent = total.toFixed(2);
        var total2 = (total - abona.value);
        resto.value = total2.toFixed(2);
        }
    }.bind(null, total));

    
}


function generarPDF(cliente, id_venta) {
    url = 'pdf/generar.php?cl=' + cliente + '&v=' + id_venta;
    window.open(url, '_blank');
}

function btnCambiar(e) {
    e.preventDefault();
    const actual = document.getElementById('actual').value;
    const nueva = document.getElementById('nueva').value;
    if (actual == "" || nueva == "") {
        Swal.fire({
            position: 'top-end',
            icon: 'error',
            title: 'Los campos estan vacios',
            showConfirmButton: false,
            timer: 2000
        })
    } else {
        const cambio = 'pass';
        $.ajax({
            url: "ajax.php",
            type: 'POST',
            data: {
                actual: actual,
                nueva: nueva,
                cambio: cambio
            },
            success: function (response) {
                console.log(response);
                if (response == 'ok') {
                    Swal.fire({
                        position: 'top-end',
                        icon: 'success',
                        title: 'Contraseña modificado',
                        showConfirmButton: false,
                        timer: 2000
                    })
                    document.querySelector('frmPass').reset();
                    $("#nuevo_pass").modal("hide");
                } else if (response == 'dif') {
                    Swal.fire({
                        position: 'top-end',
                        icon: 'error',
                        title: 'La contraseña actual incorrecta',
                        showConfirmButton: false,
                        timer: 2000
                    })
                } else {
                    Swal.fire({
                        position: 'top-end',
                        icon: 'error',
                        title: 'Error al modificar la contraseña',
                        showConfirmButton: false,
                        timer: 2000
                    })
                }
            }
        });
    }
}


